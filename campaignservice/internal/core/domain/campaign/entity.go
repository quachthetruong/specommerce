package campaign

import "time"

type Campaign struct {
	Id          int64          `json:"id" validate:"required"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Policy      map[string]any `json:"policy" validate:"required"`
	StartTime   time.Time      `json:"start_time" validate:"required"`
	EndTime     time.Time      `json:"end_time" validate:"required"`
	CreatedAt   time.Time      `json:"created_at" validate:"required"`
	UpdatedAt   time.Time      `json:"updated_at" validate:"required"`
}

type IphoneCampaignPolicy struct {
	TotalReward      int64 `json:"total_reward" validate:"required"`
	MinOrderAmount   int64 `json:"min_order_amount" validate:"required"`
	MaxTrackedOrders int64 `json:"max_tracked_orders" validate:"required"`
}
