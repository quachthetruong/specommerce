package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	orderHandler "specommerce/orderservice/internal/adapters/primary/order/handler"
)

// Route for admin user
func adminRoutes(routerGroup *gin.RouterGroup, injector do.Injector) {
	order := do.MustInvoke[orderHandler.OrderHandler](injector)

	v1OrderGroup := routerGroup.Group("/v1/orders")
	v1OrderGroup.GET("", order.GetAllOrders)
}
