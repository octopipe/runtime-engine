package main

import (
	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/controller/sharedinfra"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

func init() {

}

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		Port:                   9443,
		HealthProbeBindAddress: ":8001",
		LeaderElection:         false,
		LeaderElectionID:       "dec90f54.circlerr.io",
	})
	if err != nil {
		panic(err)
	}

	terraformProvider, err := terraform.NewTerraformProvider()
	if err != nil {
		panic(err)
	}

	pluginManager := pluginmanager.NewPluginManager(terraformProvider)

	terraformController := sharedinfra.NewController(
		logger,
		mgr.GetClient(),
		mgr.GetScheme(),
		pluginManager,
	)

	if err := terraformController.SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		panic(err)
	}

	logger.Info("start sharedInfra controller")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}