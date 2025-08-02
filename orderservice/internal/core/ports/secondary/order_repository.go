package secondary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
)

// OrderRepository defines the secondary port for order persistence
type OrderRepository interface {
	Create(ctx context.Context, order order.Order) (order.Order, error)
	GetAll(ctx context.Context) ([]order.Order, error)
}
