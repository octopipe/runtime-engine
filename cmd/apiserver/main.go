package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/infra"
	"github.com/octopipe/cloudx/internal/providerconfig"
	"github.com/octopipe/cloudx/internal/taskoutput"
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
	infraRepository := infra.NewK8sRepository(k8sClient)
	infraUseCase := infra.NewUseCase(infraRepository)

	connectionInterfaceRepository := taskoutput.NewK8sRepository(k8sClient)
	connectionInterfaceUseCase := taskoutput.NewUseCase(connectionInterfaceRepository)

	providerConfigRepository := providerconfig.NewK8sRepository(k8sClient)
	providerConfigUseCase := providerconfig.NewUseCase(providerConfigRepository)

	r = infra.NewHTTPHandler(r, infraUseCase)
	r = taskoutput.NewHTTPHandler(r, connectionInterfaceUseCase)
	r = providerconfig.NewHTTPHandler(r, providerConfigUseCase)

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": ":)",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
