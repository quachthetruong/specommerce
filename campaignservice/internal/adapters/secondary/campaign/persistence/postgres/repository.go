package postgres

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
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

func (r *campaignPersistenceRepository) Update(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error) {
	errTemplate := "campaignPersistenceRepository UpdateCampaign %w"
	campaignModel, err := FromDomainModel(campaign)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	updated, err := database.NewPostgresCrudDatabaseOperation[Campaign](r.getDbFunc).Update(ctx, campaignModel)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}

	entity, err := updated.ToDomainModel()
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	return entity, nil
}

func (r *campaignPersistenceRepository) GetCampaignByType(ctx context.Context, campaignType string) (domain.Campaign, error) {
	errTemplate := "campaignPersistenceRepository GetCampaignByType %w"
	record, err := database.NewPostgresCrudDatabaseOperation[Campaign](r.getDbFunc).Get(ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("type = ?", campaignType)
	})
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	entity, err := record.ToDomainModel()
	if err != nil {
		return domain.Campaign{}, fmt.Errorf(errTemplate, err)
	}
	return entity, err
}

func (r *campaignPersistenceRepository) GetIphoneWinner(ctx context.Context, iphoneCampaign domain.IphoneCampaign) ([]domain.IphoneWinner, error) {
	errTemplate := "campaignPersistenceRepository.GetIphoneWinner: %w"
	query := `
		with first_customers as (
			select customer_id, customer_name, min(created_at) as first_order_date,
			max(total_amount) as max_order_amount
			from orders where created_at >= ? and created_at <= ?
			group by customer_id, customer_name
			order by min(created_at)
			limit ?
		)
		select customer_id, customer_name, first_order_date, max_order_amount 
		from first_customers where max_order_amount > ?
		order by first_order_date limit ?
	`

	results := make([]domain.IphoneWinner, 0)
	rows, err := r.getDbFunc(ctx).QueryContext(ctx, query,
		iphoneCampaign.StartTime,
		iphoneCampaign.EndTime,
		iphoneCampaign.Policy.MaxTrackedOrders,
		iphoneCampaign.Policy.MinOrderAmount,
		iphoneCampaign.Policy.TotalReward)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			err = fmt.Errorf(errTemplate, err)
		}
	}()

	for rows.Next() {
		var record IphoneWinner
		err = rows.Scan(&record.CustomerId, &record.CustomerName, &record.FirstOrderTime, &record.MaxTotalOrderAmount)
		if err != nil {
			return nil, fmt.Errorf(errTemplate, err)
		}
		results = append(results, record.ToDomainModel())
	}

	return results, nil
}
