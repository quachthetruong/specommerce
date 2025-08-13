package primary

import (
	"context"
	"specommerce/campaignservice/internal/core/domain/order"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	ProcessPendingOrder(ctx context.Context, order order.Order) error
	ProcessOrderResult(ctx context.Context, order order.Order) error
	SaveSuccessOrder(ctx context.Context, order order.Order) error
}
