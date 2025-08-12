package primary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/pagination"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, order order.CreateOrderRequest) (order.Order, error)
	ProcessPaymentResponse(ctx context.Context, request payment.ProcessPaymentResponse) (order.Order, error)
	GetAllOrders(ctx context.Context) ([]order.Order, error)
	SearchOrders(ctx context.Context, filter secondary.SearchOrdersFilter) (pagination.Page[order.Order], error)
}
