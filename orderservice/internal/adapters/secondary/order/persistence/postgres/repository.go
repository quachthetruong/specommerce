package postgres

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	domain "specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/database"
)

type orderPersistenceRepository struct {
	getDbFunc database.GetDbFunc
}

func NewOrderPersistenceRepository(dbFunc database.GetDbFunc) secondary.OrderRepository {
	return &orderPersistenceRepository{
		getDbFunc: dbFunc,
	}
}

func (r *orderPersistenceRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	orders, err := database.NewPostgresCrudDatabaseOperation[Order](r.getDbFunc).FindAll(ctx)
	entities := make([]domain.Order, 0, len(orders))
	if err != nil {
		return []domain.Order{}, fmt.Errorf("orderPersistenceRepository GetAllOrders %w", err)
	}
	for _, order := range orders {
		entities = append(entities, order.ToDomainModel())
	}
	return entities, nil
}

func (r *orderPersistenceRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	created, err := database.NewPostgresCrudDatabaseOperation[Order](r.getDbFunc).Create(ctx, FromDomainModel(order))
	if err != nil {
		return domain.Order{}, fmt.Errorf("orderPersistenceRepository CreateOrder %w", err)
	}
	return created.ToDomainModel(), nil
}

func (r *orderPersistenceRepository) UpdateStatusById(ctx context.Context, id xid.ID, status domain.OrderStatus) (domain.Order, error) {
	errTemplate := "orderPersistenceRepository.UpdateStatusById: %w"
	record := Order{}
	_, err := r.getDbFunc(ctx).NewUpdate().Model((*Order)(nil)).
		Where("id = ?", id).
		Set("status = ?", status).
		Returning("*").Exec(ctx, &record)
	if err != nil {
		return domain.Order{}, fmt.Errorf(errTemplate, err)
	}
	return record.ToDomainModel(), nil
}
