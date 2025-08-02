package secondary

import (
	"context"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/domain/customer"
)

type CampaignRepository interface {
	Create(ctx context.Context, input domain.Campaign) (domain.Campaign, error)
	GetWinner(ctx context.Context, campaignID int64) ([]customer.Customer, error)
}
