package primary

import (
	"context"
	"specommerce/campaignservice/internal/core/domain/customer"
	"specommerce/campaignservice/internal/core/domain/order"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, order order.Order) (order.Order, error)
	GetWinner(ctx context.Context, campaignID int64) ([]customer.Customer, error)
}
