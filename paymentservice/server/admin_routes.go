package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	paymentHandler "specommerce/paymentservice/internal/adapters/primary/payment/handler"
)

// Route for admin user
func adminRoutes(routerGroup *gin.RouterGroup, injector do.Injector) {
	payment := do.MustInvoke[paymentHandler.PaymentHandler](injector)

	v1PaymentGroup := routerGroup.Group("/v1/payments")
	v1PaymentGroup.GET("", payment.GetAllPayments)
	v1PaymentGroup.GET("/search", payment.SearchPayments)
}
