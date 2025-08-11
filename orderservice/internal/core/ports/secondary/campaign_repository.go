package secondary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
)

type CampaignRepository interface {
	SendOrderEvent(ctx context.Context, input order.Order) error
}
