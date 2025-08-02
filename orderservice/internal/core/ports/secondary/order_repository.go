package secondary

import (
	"context"
	"github.com/rs/xid"
	"specommerce/orderservice/internal/core/domain/order"
)

// OrderRepository defines the secondary port for order persistence
type OrderRepository interface {
	Create(ctx context.Context, order order.Order) (order.Order, error)
	GetAll(ctx context.Context) ([]order.Order, error)
	UpdateStatusById(ctx context.Context, id xid.ID, status order.OrderStatus) (order.Order, error)
}
