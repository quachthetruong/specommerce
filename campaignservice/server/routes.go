package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"specommerce/campaignservice/config"
	docs "specommerce/campaignservice/docs/openapi/api/orderservice"
	"specommerce/campaignservice/pkg/environment"
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
	docs.SwaggerInfo.Title = "SP Ecommerce - Order Service API"
	docs.SwaggerInfo.Description = "Order Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8080"
	r := gin.New()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

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
