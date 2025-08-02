package postgres

import (
	"context"
	"fmt"
	domain "specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/database"
)

type orderPersistenceRepository struct {
	getDbFunc database.GetDbFunc
}

func NewOrderPersistenceRepository(dbFunc database.GetDbFunc) secondary.OrderRepository {
	return &orderPersistenceRepository{
		getDbFunc: dbFunc,
	}
}

func (r *orderPersistenceRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	created, err := database.NewPostgresCrudDatabaseOperation[Order](r.getDbFunc).Create(ctx, FromDomainModel(order))
	if err != nil {
		return domain.Order{}, fmt.Errorf("orderPersistenceRepository CreateOrder %w", err)
	}
	return created.ToDomainModel(), nil
}
