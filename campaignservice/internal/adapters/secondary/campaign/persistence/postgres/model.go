package postgres

import (
	"github.com/uptrace/bun"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"time"
)

type Campaign struct {
	bun.BaseModel `bun:"campaigns"`
	Id            int64          `bun:",skipupdate,pk"`
	Name          string         `bun:"name,notnull"`
	Description   string         `bun:"description,notnull"`
	Policy        map[string]any `bun:"type:jsonb,default:'{}'::jsonb"`
	StartTime     time.Time      `bun:"start_time,notnull"`
	EndTime       time.Time      `bun:"end_time,notnull"`
	CreatedAt     time.Time      `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt     time.Time      `bun:",nullzero,notnull,default:current_timestamp"`
}

func (c Campaign) ToDomainModel() (domain.Campaign, error) {
	campaign := domain.Campaign{
		Id:          c.Id,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Name:        c.Name,
		Description: c.Description,
		StartTime:   c.StartTime,
		EndTime:     c.EndTime,
		Policy:      c.Policy,
	}

	return campaign, nil
}

func FromDomainModel(dm domain.Campaign) (Campaign, error) {
	return Campaign{
		Id:          dm.Id,
		Name:        dm.Name,
		Description: dm.Description,
		StartTime:   dm.StartTime,
		EndTime:     dm.EndTime,
		Policy:      dm.Policy,
		CreatedAt:   dm.CreatedAt,
		UpdatedAt:   dm.UpdatedAt,
	}, nil
}
