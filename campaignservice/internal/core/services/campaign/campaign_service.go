package campaign

import (
	"context"
	"fmt"
	"specommerce/campaignservice/config"
	"specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/atomicity"
	"specommerce/campaignservice/pkg/cache"
	"strconv"
)

type campaignService struct {
	campaignRepository secondary.CampaignRepository
	atomicExecutor     atomicity.AtomicExecutor
	config             config.AppConfig
	cacheClient        cache.Cache
}

func NewCampaignService(
	campaignRepository secondary.CampaignRepository,
	atomicExecutor atomicity.AtomicExecutor,
	config config.AppConfig,
	cacheClient cache.Cache,
) primary.CampaignService {
	return &campaignService{
		campaignRepository: campaignRepository,
		atomicExecutor:     atomicExecutor,
		config:             config,
		cacheClient:        cacheClient,
	}
}

func (s *campaignService) CreateCampaign(ctx context.Context, input campaign.Campaign) (campaign.Campaign, error) {
	errTemplate := "campaignService Create %w"
	savedCampaign, err := s.campaignRepository.Create(ctx, input)
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf(errTemplate, err)
	}

	// Store full campaign information in Redis Hash using Lua script
	luaScript := `
		local key = KEYS[1]
		local id = ARGV[1]
		local name = ARGV[2] 
		local type = ARGV[3]
		local description = ARGV[4]
		local total_reward = ARGV[5]
		local min_order_amount = ARGV[6]
		local max_tracked_orders = ARGV[7]
		local start_time_micro = ARGV[8]
		local end_time_micro = ARGV[9]
		local created_at_micro = ARGV[10]
		local updated_at_micro = ARGV[11]
		
		return redis.call('HMSET', key,
			'id', id,
			'name', name,
			'type', type,
			'description', description,
			'policy_total_reward', total_reward,
			'policy_min_order_amount', min_order_amount,
			'policy_max_tracked_orders', max_tracked_orders,
			'start_time_micro', start_time_micro,
			'end_time_micro', end_time_micro,
			'created_at_micro', created_at_micro,
			'updated_at_micro', updated_at_micro
		)
	`

	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)

	// Extract policy fields from the map
	totalReward := fmt.Sprintf("%.0f", savedCampaign.Policy["total_reward"])
	minOrderAmount := fmt.Sprintf("%.0f", savedCampaign.Policy["min_order_amount"])
	maxTrackedOrders := fmt.Sprintf("%.0f", savedCampaign.Policy["max_tracked_orders"])

	_, err = s.cacheClient.Eval(ctx, luaScript, []string{campaignKey},
		strconv.FormatInt(savedCampaign.Id, 10),
		savedCampaign.Name,
		savedCampaign.Type,
		savedCampaign.Description,
		totalReward,
		minOrderAmount,
		maxTrackedOrders,
		strconv.FormatInt(savedCampaign.StartTime.UnixMicro(), 10),
		strconv.FormatInt(savedCampaign.EndTime.UnixMicro(), 10),
		strconv.FormatInt(savedCampaign.CreatedAt.UnixMicro(), 10),
		strconv.FormatInt(savedCampaign.UpdatedAt.UnixMicro(), 10),
	)
	if err != nil {
		// Log error but don't fail the campaign creation
		fmt.Printf("Failed to store campaign in Redis: %v\n", err)
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

	// Update campaign information in Redis Hash using Lua script
	luaScript := `
		local key = KEYS[1]
		local id = ARGV[1]
		local name = ARGV[2] 
		local type = ARGV[3]
		local description = ARGV[4]
		local total_reward = ARGV[5]
		local min_order_amount = ARGV[6]
		local max_tracked_orders = ARGV[7]
		local start_time_micro = ARGV[8]
		local end_time_micro = ARGV[9]
		local created_at_micro = ARGV[10]
		local updated_at_micro = ARGV[11]
		
		return redis.call('HMSET', key,
			'id', id,
			'name', name,
			'type', type,
			'description', description,
			'policy_total_reward', total_reward,
			'policy_min_order_amount', min_order_amount,
			'policy_max_tracked_orders', max_tracked_orders,
			'start_time_micro', start_time_micro,
			'end_time_micro', end_time_micro,
			'created_at_micro', created_at_micro,
			'updated_at_micro', updated_at_micro
		)
	`
	
	campaignKey := fmt.Sprintf("campaign:%s", s.config.IphoneCampaign)
	
	// Extract policy fields from the map
	totalReward := fmt.Sprintf("%.0f", updatedCampaign.Policy["total_reward"])
	minOrderAmount := fmt.Sprintf("%.0f", updatedCampaign.Policy["min_order_amount"])
	maxTrackedOrders := fmt.Sprintf("%.0f", updatedCampaign.Policy["max_tracked_orders"])
	
	_, err = s.cacheClient.Eval(ctx, luaScript, []string{campaignKey}, 
		strconv.FormatInt(updatedCampaign.Id, 10),
		updatedCampaign.Name,
		updatedCampaign.Type, 
		updatedCampaign.Description,
		totalReward,
		minOrderAmount,
		maxTrackedOrders,
		strconv.FormatInt(updatedCampaign.StartTime.UnixMicro(), 10),
		strconv.FormatInt(updatedCampaign.EndTime.UnixMicro(), 10),
		strconv.FormatInt(updatedCampaign.CreatedAt.UnixMicro(), 10),
		strconv.FormatInt(updatedCampaign.UpdatedAt.UnixMicro(), 10),
	)
	if err != nil {
		// Log error but don't fail the campaign update
		fmt.Printf("Failed to update campaign in Redis: %v\n", err)
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
