package primary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, order order.Order) (order.Order, error)
	GetAllOrders(ctx context.Context) ([]order.Order, error)
}
