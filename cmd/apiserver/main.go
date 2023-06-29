package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/connectioninterface"
	"github.com/octopipe/cloudx/internal/execution"
	"github.com/octopipe/cloudx/internal/providerconfig"
	"github.com/octopipe/cloudx/internal/sharedinfra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

func main() {
	// logger, _ := zap.NewProduction()
	_ = godotenv.Load()

	config := ctrl.GetConfigOrDie()
	k8sClient, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(CORSMiddleware())
	executionRepository := execution.NewK8sRepository(k8sClient)
	executionUseCase := execution.NewUseCase(executionRepository)

	sharedInfraRepository := sharedinfra.NewK8sRepository(k8sClient)
	sharedInfraUseCase := sharedinfra.NewUseCase(sharedInfraRepository)

	connectionInterfaceRepository := connectioninterface.NewK8sRepository(k8sClient)
	connectionInterfaceUseCase := connectioninterface.NewUseCase(connectionInterfaceRepository)

	providerConfigRepository := providerconfig.NewK8sRepository(k8sClient)
	providerConfigUseCase := providerconfig.NewUseCase(providerConfigRepository)

	r = execution.NewHTTPHandler(r, executionUseCase)
	r = sharedinfra.NewHTTPHandler(r, sharedInfraUseCase)
	r = connectioninterface.NewHTTPHandler(r, connectionInterfaceUseCase)
	r = providerconfig.NewHTTPHandler(r, providerConfigUseCase)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
