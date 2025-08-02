package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"specommerce/paymentservice/config"
	docs "specommerce/paymentservice/docs/openapi/api/paymentservice"
	"specommerce/paymentservice/pkg/environment"
)

// NewRoutes godoc
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func routes(injector do.Injector) http.Handler {
	appConfig := do.MustInvoke[config.AppConfig](injector)
	if appConfig.Env == string(environment.Production) {
		gin.SetMode(gin.ReleaseMode)
	}
	docs.SwaggerInfo.Title = "SP Ecommerce - Payment Service API"
	docs.SwaggerInfo.Description = "Payment Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8081"
	r := gin.New()

	r.NoRoute(notFound)
	r.NoMethod(methodNotAllowed)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(
		"/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)
	apiUserGroup := r.Group("/api")
	consumerRoutes(apiUserGroup, injector)

	apiAdminGroup := r.Group("/api/admin")
	adminRoutes(apiAdminGroup, injector)
	return r
}
