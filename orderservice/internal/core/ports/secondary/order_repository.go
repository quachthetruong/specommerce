package secondary

import (
	"context"
	"github.com/rs/xid"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/pkg/pagination"
)

// SearchOrdersFilter represents the filter for searching orders
type SearchOrdersFilter struct {
	pagination.Paging
}

// OrderRepository defines the secondary port for order persistence
type OrderRepository interface {
	Create(ctx context.Context, order order.Order) (order.Order, error)
	GetAll(ctx context.Context) ([]order.Order, error)
	UpdateStatusById(ctx context.Context, id xid.ID, status order.OrderStatus) (order.Order, error)
	SearchOrders(ctx context.Context, filter SearchOrdersFilter) (pagination.Page[order.Order], error)
}
