package engine

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/connectioninterface"
	"github.com/octopipe/cloudx/internal/plugin"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ThisInterpolationOrigin                = "this"
	ConnectionInterfaceInterpolationOrigin = "connection-interface"
)

type DependencyGraph map[string][]string

type ExecutionContextItem struct {
	Value     string
	Sensitive bool
	Type      string
}

type pipeline struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex
	executionContext  map[string]map[string]ExecutionContextItem
	rpcClient         rpcclient.Client
}

func NewPipeline(logger *zap.Logger, rpcClient rpcclient.Client, terraformProvider terraform.TerraformProvider) pipeline {
	return pipeline{
		logger:            logger,
		terraformProvider: terraformProvider,
		executionContext:  make(map[string]map[string]ExecutionContextItem),
		rpcClient:         rpcClient,
	}
}

func (p *pipeline) Execute(action ExecutionActionType, graph DependencyGraph, sharedInfra commonv1alpha1.SharedInfra, currentExecutionStatusChann chan commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
	lastExecution := sharedInfra.Status.LastExecution
	status := commonv1alpha1.ExecutionStatus{
		Status:  ExecutionRunningStatus,
		Plugins: []commonv1alpha1.PluginExecutionStatus{},
	}

	eg := new(errgroup.Group)
	inDegrees := make(map[string]int)

	for node, deps := range graph {
		inDegrees[node] = len(deps)
	}

	for {
		for node, deps := range inDegrees {
			if _, ok := p.executionContext[node]; !ok && deps == 0 {
				eg.Go(func(node string) func() error {
					return func() error {
						p.logger.Info("starting plugin execution...", zap.String("name", node), zap.Any("action", action))

						pluginExecutionStatus, pluginOutput := commonv1alpha1.PluginExecutionStatus{}, map[string]ExecutionContextItem{}
						if action == DestroyAction {
							lastPluginExecution := commonv1alpha1.PluginExecutionStatus{}
							for _, statusPlugin := range lastExecution.Plugins {
								if statusPlugin.Name == node {
									lastPluginExecution = statusPlugin
									break
								}
							}
							pluginExecutionStatus = p.destroyPlugin(lastExecution, lastPluginExecution)
						} else {
							currentPlugin := commonv1alpha1.SharedInfraPlugin{}
							for _, specPlugin := range sharedInfra.Spec.Plugins {
								if specPlugin.Name == node {
									currentPlugin = specPlugin
									break
								}
							}
							pluginExecutionStatus, pluginOutput = p.applyPlugin(lastExecution, currentPlugin)
						}

						p.mu.Lock()
						defer p.mu.Unlock()

						status.Plugins = append(status.Plugins, pluginExecutionStatus)
						if pluginExecutionStatus.Status == ExecutionApplyErrorStatus || pluginExecutionStatus.Status == ExecutionDestroyErrorStatus {
							status.Status = ExecutionErrorStatus
							status.Error = pluginExecutionStatus.Error
							p.logger.Info("plugin execution failed", zap.String("plugin-name", pluginExecutionStatus.Name), zap.Error(errors.New(pluginExecutionStatus.Error)))
							return errors.New(pluginExecutionStatus.Error)
						}

						p.executionContext[node] = pluginOutput

						for n, deps := range graph {
							for _, dep := range deps {
								if dep == node {
									inDegrees[n]--
								}
							}
						}

						p.logger.Info("finish plugin execution...", zap.String("name", node), zap.Any("action", action))
						return nil
					}
				}(node))

			}
		}

		err := eg.Wait()
		currentExecutionStatusChann <- status
		if err != nil {
			p.logger.Info("find errors in parallel execution...")
			break
		}

		if len(p.executionContext) == len(graph) {
			break
		}
	}
	status.Status = ExecutionSuccessStatus
	p.logger.Info("finished pipeline execution")
	return status
}

func (p *pipeline) destroyPlugin(lastExecution commonv1alpha1.ExecutionStatus, lastExecutionPlugin commonv1alpha1.PluginExecutionStatus) commonv1alpha1.PluginExecutionStatus {
	status := commonv1alpha1.PluginExecutionStatus{
		Name:       lastExecutionPlugin.Name,
		Ref:        lastExecutionPlugin.Ref,
		Depends:    lastExecutionPlugin.Depends,
		Inputs:     lastExecutionPlugin.Inputs,
		PluginType: lastExecutionPlugin.PluginType,
		Status:     ExecutionDestroyed,
		StartedAt:  time.Now().Format(time.RFC3339),
	}

	inputs := lastExecutionPlugin.Inputs
	if lastExecutionPlugin.PluginType == plugin.TerraformPluginType {
		err := p.terraformProvider.Destroy(lastExecutionPlugin.Ref, inputs, lastExecutionPlugin.State, lastExecutionPlugin.DependencyLock)
		if err != nil {
			status.Error = err.Error()
			status.Status = ExecutionDestroyErrorStatus
			return status
		}

		return status
	}

	status.Error = "invalid plugin type"
	status.Status = ExecutionDestroyErrorStatus

	return status
}

func (p *pipeline) applyPlugin(lastExecution commonv1alpha1.ExecutionStatus, currentPlugin commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginExecutionStatus, map[string]ExecutionContextItem) {
	lastPluginExecutionStatus := commonv1alpha1.PluginExecutionStatus{}

	for _, e := range lastExecution.Plugins {
		if e.Name == currentPlugin.Name {
			lastPluginExecutionStatus = e
		}
	}

	status := commonv1alpha1.PluginExecutionStatus{
		Name:       currentPlugin.Name,
		Ref:        currentPlugin.Ref,
		Depends:    currentPlugin.Depends,
		Inputs:     currentPlugin.Inputs,
		PluginType: currentPlugin.PluginType,
		Status:     ExecutionAppliedStatus,
		StartedAt:  time.Now().Format(time.RFC3339),
	}

	inputs, err := p.interpolatePluginInputsByExecutionContext(currentPlugin)
	if err != nil {
		status.Error = err.Error()
		status.Status = ExecutionApplyErrorStatus
		return status, nil
	}

	status.Inputs = inputs
	if currentPlugin.PluginType == plugin.TerraformPluginType {
		out, state, lockfile, err := p.terraformProvider.Apply(currentPlugin.Ref, inputs, lastPluginExecutionStatus.State, lastPluginExecutionStatus.DependencyLock)
		status.FinishedAt = time.Now().Format(time.RFC3339)
		if err != nil {
			status.Error = err.Error()
			status.Status = ExecutionApplyErrorStatus
			return status, nil
		}

		status.DependencyLock = lockfile
		status.State = state

		outputs := map[string]ExecutionContextItem{}

		for key, tfMeta := range out {
			outputs[key] = ExecutionContextItem{
				Value:     string(tfMeta.Value),
				Type:      string(tfMeta.Type),
				Sensitive: tfMeta.Sensitive,
			}
		}

		return status, outputs
	}

	status.Error = "invalid plugin type"
	status.Status = ExecutionApplyErrorStatus

	return status, nil
}

func (p *pipeline) interpolatePluginInputsByExecutionContext(plugin commonv1alpha1.SharedInfraPlugin) ([]commonv1alpha1.SharedInfraPluginInput, error) {
	inputs := []commonv1alpha1.SharedInfraPluginInput{}
	for _, i := range plugin.Inputs {
		tokens := Lex(i.Value)
		data := map[string]string{}
		sensitive := false
		for _, t := range tokens {
			if t.Type == TokenVariable {
				s := strings.Split(strings.Trim(t.Value, " "), ".")
				if len(s) != 3 {
					return nil, fmt.Errorf("malformed input variable %s with value %s", i.Key, i.Value)
				}

				value, isSensitive, err := p.getDataByOrigin(s[0], s[1], s[2])
				if err != nil {
					return nil, err
				}

				if isSensitive {
					sensitive = isSensitive
				}

				data[t.Value] = strings.Trim(value, "\"")
			}
		}

		inputs = append(inputs, commonv1alpha1.SharedInfraPluginInput{
			Key:       i.Key,
			Value:     Interpolate(tokens, data),
			Sensitive: sensitive,
		})
	}

	return inputs, nil
}

func (p *pipeline) getDataByOrigin(origin string, name string, attr string) (string, bool, error) {
	switch origin {
	case ThisInterpolationOrigin:
		p.logger.Info("interpolate this origin")
		execution, ok := p.executionContext[name]
		if !ok {
			return "", false, fmt.Errorf("not found plugin %s in execution context", name)
		}

		executionAttr, ok := execution[attr]
		if !ok {
			return "", false, fmt.Errorf("not found attr %s in finished plugin execution %s", attr, name)
		}

		return executionAttr.Value, executionAttr.Sensitive, nil

	case ConnectionInterfaceInterpolationOrigin:
		p.logger.Info("interpolate this connection-interface")
		connectionInterface := commonv1alpha1.ConnectionInterface{}
		err := p.rpcClient.Call("ConnectionInterfaceRPCHandler.GetConnectionInterface", connectioninterface.RPCGetConnectionInterfaceArgs{
			Ref: types.NamespacedName{Name: name, Namespace: "default"},
		}, &connectionInterface)
		if err != nil {
			return "", false, err
		}

		for _, out := range connectionInterface.Spec.Outputs {
			if out.Key == attr {
				return out.Value, out.Sensitive, nil
			}
		}

		return "", false, fmt.Errorf("not found attr in connection-interface %s", name)
	default:
		return "", false, fmt.Errorf("invalid origin type %s", origin)
	}
}
