package secondary

import (
	"context"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
)

type CampaignRepository interface {
	Create(ctx context.Context, input domain.Campaign) (domain.Campaign, error)
	Update(ctx context.Context, input domain.Campaign) (domain.Campaign, error)
	GetIphoneWinner(ctx context.Context, campaign domain.IphoneCampaign) ([]domain.IphoneWinner, error)
	GetCampaignByType(ctx context.Context, campaignType string) (domain.Campaign, error)
	SaveWinner(ctx context.Context, campaignId int64, customerId string) error
}
