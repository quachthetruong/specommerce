package campaign

import (
	"context"
	"fmt"
	"specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/domain/customer"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/atomicity"
)

type campaignService struct {
	campaignRepository secondary.CampaignRepository
	atomicExecutor     atomicity.AtomicExecutor
}

func NewCampaignService(
	campaignRepository secondary.CampaignRepository,
	atomicExecutor atomicity.AtomicExecutor,
) primary.CampaignService {
	return &campaignService{
		campaignRepository: campaignRepository,
		atomicExecutor:     atomicExecutor,
	}
}

func (s *campaignService) CreateCampaign(ctx context.Context, input campaign.Campaign) (campaign.Campaign, error) {
	errTemplate := "campaignService Create %w"
	savedCampaign, err := s.campaignRepository.Create(ctx, input)
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	return savedCampaign, nil
}

func (s *campaignService) GetWinner(ctx context.Context, campaignID int64) ([]customer.Customer, error) {
	errTemplate := "campaignService GetWinner %w"
	winners, err := s.campaignRepository.GetWinner(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	return winners, nil
}
