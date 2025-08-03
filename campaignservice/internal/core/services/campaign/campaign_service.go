package campaign

import (
	"context"
	"fmt"
	"specommerce/campaignservice/config"
	"specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/atomicity"
)

type campaignService struct {
	campaignRepository secondary.CampaignRepository
	atomicExecutor     atomicity.AtomicExecutor
	config             config.AppConfig
}

func NewCampaignService(
	campaignRepository secondary.CampaignRepository,
	atomicExecutor atomicity.AtomicExecutor,
	config config.AppConfig,
) primary.CampaignService {
	return &campaignService{
		campaignRepository: campaignRepository,
		atomicExecutor:     atomicExecutor,
		config:             config,
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

func (s *campaignService) GetIphoneCampaign(ctx context.Context) (campaign.Campaign, error) {
	errTemplate := "campaignService GetIphoneCampaign %w"
	campaign, err := s.campaignRepository.GetCampaignByType(ctx, s.config.IphoneCampaign)
	if err != nil {
		return campaign, fmt.Errorf(errTemplate, err)
	}
	return campaign, nil
}

func (s *campaignService) UpdateIphoneCampaign(ctx context.Context, input campaign.Campaign) (campaign.Campaign, error) {
	errTemplate := "campaignService UpdateIphoneCampaign %w"
	updatedCampaign, err := s.campaignRepository.Update(ctx, input)
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	return updatedCampaign, nil
}

func (s *campaignService) GetIphoneWinner(ctx context.Context) ([]campaign.IphoneWinner, error) {
	errTemplate := "campaignService GetWinner %w"
	campaign, err := s.campaignRepository.GetCampaignByType(ctx, s.config.IphoneCampaign)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	iphoneCampaign, err := campaign.ToIphoneCampaign()
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	winners, err := s.campaignRepository.GetIphoneWinner(ctx, iphoneCampaign)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	return winners, nil
}
