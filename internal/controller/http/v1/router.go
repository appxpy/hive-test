// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/appxpy/hive-test/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/appxpy/hive-test/docs"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/appxpy/hive-test/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Hive-Test
// @description Тестовое задание для вступления в проект
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(
	handler *gin.Engine,
	config *config.Config,
	l logger.Interface,
	u usecase.UserUseCase,
	a usecase.AssetUseCase,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, u, l)
		newAssetRoutes(h, a, l, config.JWTSecret)
	}
}
