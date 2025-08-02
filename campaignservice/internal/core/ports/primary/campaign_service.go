package primary

import (
	"context"
	"specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/domain/customer"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, input campaign.Campaign) (campaign.Campaign, error)
	GetWinner(ctx context.Context, campaignID int64) ([]customer.Customer, error)
}
