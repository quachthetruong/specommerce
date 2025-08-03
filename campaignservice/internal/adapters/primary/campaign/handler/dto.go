package handler

import (
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"time"
)

// CreateCampaignRequest represents the request for creating a campaign
type CreateIphoneCampaignRequest struct {
	Name             string    `json:"name" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	StartTime        time.Time `json:"start_time" binding:"required"`
	EndTime          time.Time `json:"end_time" binding:"required"`
	TotalReward      int64     `json:"total_reward" binding:"required"`
	MinOrderAmount   int64     `json:"min_order_amount" binding:"required"`
	MaxTrackedOrders int64     `json:"max_tracked_orders" binding:"required"`
}

// UpdateCampaignRequest represents the request for updating a campaign
type UpdateIphoneCampaignRequest struct {
	Id               int64     `uri:"id" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	StartTime        time.Time `json:"start_time" binding:"required"`
	EndTime          time.Time `json:"end_time" binding:"required"`
	TotalReward      int64     `json:"total_reward" binding:"required"`
	MinOrderAmount   int64     `json:"min_order_amount" binding:"required"`
	MaxTrackedOrders int64     `json:"max_tracked_orders" binding:"required"`
}

func (r CreateIphoneCampaignRequest) ToDomain(campaignType string) domain.Campaign {
	return domain.Campaign{
		Name:        r.Name,
		Type:        campaignType,
		Description: r.Description,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Policy: map[string]any{
			"total_reward":       r.TotalReward,
			"min_order_amount":   r.MinOrderAmount,
			"max_tracked_orders": r.MaxTrackedOrders,
		},
	}
}

func (r UpdateIphoneCampaignRequest) ToDomain(id int64) domain.Campaign {
	return domain.Campaign{
		Id:          id,
		Name:        r.Name,
		Type:        "iphone",
		Description: r.Description,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		UpdatedAt:   time.Now(),
		Policy: map[string]any{
			"total_reward":       r.TotalReward,
			"min_order_amount":   r.MinOrderAmount,
			"max_tracked_orders": r.MaxTrackedOrders,
		},
	}
}
