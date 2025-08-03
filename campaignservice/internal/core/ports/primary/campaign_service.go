package primary

import (
	"context"
	"specommerce/campaignservice/internal/core/domain/campaign"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, input campaign.Campaign) (campaign.Campaign, error)
	GetIphoneWinner(ctx context.Context) ([]campaign.IphoneWinner, error)
}
