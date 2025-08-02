package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"time"
)

type Campaign struct {
	bun.BaseModel `bun:"campaigns"`
	Id            int64     `bun:",skipupdate,pk"`
	Name          string    `bun:"name,notnull"`
	Description   string    `bun:"description,notnull"`
	Policy        string    `bun:"type:jsonb,default:'{}'::jsonb"`
	StartTime     time.Time `bun:"start_time,notnull"`
	EndTime       time.Time `bun:"end_time,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (c Campaign) ToDomainModel() (domain.Campaign, error) {
	errTemplate := "Campaign.ToDomainModel: %w"
	var policy domain.CampaignPolicy
	err := json.Unmarshal([]byte(c.Policy), &policy)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	campaign := domain.Campaign{
		Id:          c.Id,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Name:        c.Name,
		Description: c.Description,
		StartTime:   c.StartTime,
		EndTime:     c.EndTime,
		Policy:      policy,
	}

	return campaign, nil
}

func FromDomainModel(dm domain.Campaign) (Campaign, error) {
	errTemplate := "FromDomainModel: %w"
	policy, err := json.Marshal(dm.Policy)
	if err != nil {
		return Campaign{}, fmt.Errorf(errTemplate, err)
	}

	return Campaign{
		Id:          dm.Id,
		Name:        dm.Name,
		Description: dm.Description,
		StartTime:   dm.StartTime,
		EndTime:     dm.EndTime,
		Policy:      string(policy),
		CreatedAt:   dm.CreatedAt,
		UpdatedAt:   dm.UpdatedAt,
	}, nil
}
