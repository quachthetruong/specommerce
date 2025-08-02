package postgres

import (
	"context"
	"fmt"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/domain/customer"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/database"
)

type campaignPersistenceRepository struct {
	getDbFunc database.GetDbFunc
}

func NewCampaignPersistenceRepository(dbFunc database.GetDbFunc) secondary.CampaignRepository {
	return &campaignPersistenceRepository{
		getDbFunc: dbFunc,
	}
}

func (r *campaignPersistenceRepository) Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error) {
	errTemplate := "campaignPersistenceRepository CreateCampaign %w"
	campaignModel, err := FromDomainModel(campaign)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	created, err := database.NewPostgresCrudDatabaseOperation[Campaign](r.getDbFunc).Create(ctx, campaignModel)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}

	entity, err := created.ToDomainModel()
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	return entity, nil
}

func (r *campaignPersistenceRepository) GetWinner(ctx context.Context, campaignId int64) ([]customer.Customer, error) {
	//errTemplate := "campaignPersistenceRepository GetWinner %w"
	return []customer.Customer{}, nil
}
