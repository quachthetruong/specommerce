package campaign

import (
	"encoding/json"
	"time"
)

type Campaign struct {
	Id          int64          `json:"id" validate:"required"`
	Name        string         `json:"name" validate:"required"`
	Type        string         `json:"type" validate:"required,oneof=iphone"` // Currently only supports "iphone"
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

func (c Campaign) ToIphoneCampaign() (IphoneCampaign, error) {
	var policy IphoneCampaignPolicy
	policyBytes, err := json.Marshal(c.Policy)
	if err != nil {
		return IphoneCampaign{}, err
	}
	err = json.Unmarshal(policyBytes, &policy)
	if err != nil {
		return IphoneCampaign{}, err
	}

	return IphoneCampaign{
		Id:          c.Id,
		Name:        c.Name,
		Type:        c.Type,
		Description: c.Description,
		Policy:      policy,
		StartTime:   c.StartTime,
		EndTime:     c.EndTime,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}

type IphoneCampaign struct {
	Id          int64
	Name        string
	Type        string
	Description string
	Policy      IphoneCampaignPolicy
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type IphoneWinner struct {
	CustomerId          string    `json:"customer_id" validate:"required"`
	CustomerName        string    `json:"customer_name" validate:"required"`
	FirstOrderTime      time.Time `json:"first_order_time" validate:"required"`
	MaxTotalOrderAmount float64   `json:"max_total_order_amount" validate:"required"`
}
