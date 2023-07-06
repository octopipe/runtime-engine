package execution

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Runner struct {
	Pod     *v1.Pod
	Service *v1.Service
}

func (c *controller) NewRunner(execution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra, rawSharedInfra string, providerConfig commonv1alpha1.ProviderConfig) (Runner, error) {
	vFalse := false
	vTrue := true
	vUser := int64(65532)
	securityContext := &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Drop: []v1.Capability{"ALL"},
		},
		AllowPrivilegeEscalation: &vFalse,
		RunAsNonRoot:             &vTrue,
		RunAsUser:                &vUser,
		SeccompProfile: &v1.SeccompProfile{
			Type: v1.SeccompProfileTypeRuntimeDefault,
		},
		ReadOnlyRootFilesystem: &vTrue,
	}

	podVolumes := []v1.Volume{
		{
			Name: "temp",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "home",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}

	podVolumeMounts := []v1.VolumeMount{
		{
			Name:      "temp",
			MountPath: "/tmp",
		},
		{
			Name:      "home",
			MountPath: "/home/runner",
		},
	}
	serviceAccount := "controller-sa"

	if sharedInfra.Spec.RunnerConfig.ServiceAccount != "" {
		serviceAccount = sharedInfra.Spec.RunnerConfig.ServiceAccount
	}

	varsCreds, err := c.getCreds(providerConfig)
	if err != nil {
		return Runner{}, err
	}

	defaultVars := []v1.EnvVar{
		{
			Name:  "TF_VERSION",
			Value: "latest",
		},
		{
			Name:  "RPC_SERVER_ADDRESS",
			Value: os.Getenv("RPC_SERVER_ADDRESS"),
		},
	}

	defaultVars = append(defaultVars, varsCreds...)

	executionRef := types.NamespacedName{
		Name:      execution.Name,
		Namespace: execution.Namespace,
	}

	args := []string{execution.Spec.Action, executionRef.String(), rawSharedInfra}

	newRunnerObject := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-runner-%d", sharedInfra.GetName(), time.Now().Unix()),
			Namespace: "default",
			Labels: map[string]string{
				"commons.cloudx.io/sharedinfra-name":      sharedInfra.GetName(),
				"commons.cloudx.io/sharedinfra-namespace": sharedInfra.GetNamespace(),
				"app.kubernetes.io/managed-by":            "cloudx",
			},
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				{
					Name:            "runner",
					Image:           "mayconjrpacheco/cloudx-runner:latest",
					Args:            args,
					ImagePullPolicy: v1.PullAlways,
					SecurityContext: securityContext,
					Env:             defaultVars,
					VolumeMounts:    podVolumeMounts,
				},
			},
			Volumes: podVolumes,
		},
	}

	if os.Getenv("ENV") != "local" {
		newRunnerObject.Namespace = "cloudx-system"
		newRunnerObject.Spec.ServiceAccountName = serviceAccount
	}

	return Runner{
		Pod: newRunnerObject,
	}, nil
}

func (c controller) getCreds(providerConfig commonv1alpha1.ProviderConfig) ([]v1.EnvVar, error) {
	if providerConfig.Spec.Type == "AWS" {
		creds, err := c.provider.GetCreds(context.Background(), providerConfig)
		if err != nil {
			return nil, err
		}

		vars := []v1.EnvVar{
			{Name: "AWS_ACCESS_KEY_ID", Value: creds.AccessKeyId},
			{Name: "AWS_SECRET_ACCESS_KEY", Value: creds.AccessKey},
			{Name: "AWS_SESSION_TOKEN", Value: creds.SessionToken},
		}

		return vars, nil
	}

	return nil, errors.New("invalid provider type")
}