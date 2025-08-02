package primary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, order order.Order) (order.Order, error)
	ProcessPaymentResponse(ctx context.Context, request payment.ProcessPaymentResponse) (order.Order, error)
	GetAllOrders(ctx context.Context) ([]order.Order, error)
}
