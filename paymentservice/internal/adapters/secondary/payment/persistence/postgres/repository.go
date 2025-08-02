package postgres

import (
	"context"
	"fmt"
	domain "specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/internal/core/ports/secondary"
	"specommerce/paymentservice/pkg/database"
)

type paymentPersistenceRepository struct {
	getDbFunc database.GetDbFunc
}

func NewPaymentPersistenceRepository(dbFunc database.GetDbFunc) secondary.PaymentRepository {
	return &paymentPersistenceRepository{
		getDbFunc: dbFunc,
	}
}

func (r *paymentPersistenceRepository) GetAll(ctx context.Context) ([]domain.Payment, error) {
	payments, err := database.NewPostgresCrudDatabaseOperation[Payment](r.getDbFunc).FindAll(ctx)
	entities := make([]domain.Payment, 0, len(payments))
	if err != nil {
		return []domain.Payment{}, fmt.Errorf("paymentPersistenceRepository GetAllPayments %w", err)
	}
	for _, payment := range payments {
		entities = append(entities, payment.ToDomainModel())
	}
	return entities, nil
}

func (r *paymentPersistenceRepository) Create(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	created, err := database.NewPostgresCrudDatabaseOperation[Payment](r.getDbFunc).Create(ctx, FromDomainModel(payment))
	if err != nil {
		return domain.Payment{}, fmt.Errorf("paymentPersistenceRepository CreatePayment %w", err)
	}
	return created.ToDomainModel(), nil
}
