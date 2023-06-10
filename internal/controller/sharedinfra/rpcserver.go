package sharedinfra

import (
	"context"
	"fmt"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RPCServer struct {
	client.Client
	logger *zap.Logger
}

func NewRPCServer(client client.Client, logger *zap.Logger) *RPCServer {
	return &RPCServer{Client: client, logger: logger}
}

type RPCSetRunnerFinishedArgs struct {
	Ref         types.NamespacedName
	ExecutionId string
	Plugins     []commonv1alpha1.PluginStatus `json:"plugins,omitempty"`
	FinishedAt  string                        `json:"finishedAt,omitempty"`
	Status      string                        `json:"status,omitempty"`
	Error       string                        `json:"error,omitempty"`
}

func (s *RPCServer) SetRunnerFinished(args *RPCSetRunnerFinishedArgs, reply *int) error {
	s.logger.Info("Received rpc call", zap.String("sharedinfra", args.Ref.String()))
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, currentSharedInfra)
	if err != nil {
		return err
	}

	s.logger.Info("rpc execution", zap.String("status", args.Status))

	newExecutions := []commonv1alpha1.SharedInfraExecutionStatus{}
	allExecutions := currentSharedInfra.Status.Executions

	currentExecution := commonv1alpha1.SharedInfraExecutionStatus{}

	for _, e := range allExecutions {
		if e.Id == args.ExecutionId {
			currentExecution = commonv1alpha1.SharedInfraExecutionStatus{
				Id:         e.Id,
				StartedAt:  e.StartedAt,
				Status:     args.Status,
				FinishedAt: args.FinishedAt,
				Error:      args.Error,
				Plugins:    args.Plugins,
			}
		} else {
			newExecutions = append(newExecutions, e)
		}
	}

	newExecutions = append(newExecutions, currentExecution)
	currentSharedInfra.Status.Executions = newExecutions

	return s.Status().Update(context.TODO(), currentSharedInfra)
}

type RPCSetRunnerTimeoutArgs struct {
	SharedInfraRef types.NamespacedName
	ExecutionId    string
	FinishedAt     string
}

func (s *RPCServer) SetRunnerTimeout(args *RPCSetRunnerTimeoutArgs, reply *int) error {
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.SharedInfraRef, currentSharedInfra)
	if err != nil {
		return err
	}

	runnerList := &v1.PodList{}
	selector, _ := labels.Parse(fmt.Sprintf("commons.cloudx.io/execution-id=%s", args.ExecutionId))
	err = s.List(context.Background(), runnerList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return err
	}

	for _, r := range runnerList.Items {
		err = s.Delete(context.Background(), &r)
		if err != nil {
			return err
		}
	}

	allExecutions := currentSharedInfra.Status.Executions
	currentExecution := commonv1alpha1.SharedInfraExecutionStatus{}
	for _, e := range allExecutions {
		if e.Id == args.ExecutionId {
			currentExecution = commonv1alpha1.SharedInfraExecutionStatus{
				Id:         e.Id,
				StartedAt:  e.StartedAt,
				Status:     "TIMEOUT",
				FinishedAt: args.FinishedAt,
				Error:      "runner time exceeded",
			}
		}
	}
	allExecutions = append(allExecutions, currentExecution)
	currentSharedInfra.Status.Executions = allExecutions

	return s.Status().Update(context.TODO(), currentSharedInfra)
}
