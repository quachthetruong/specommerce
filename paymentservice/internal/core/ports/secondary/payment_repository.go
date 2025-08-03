package secondary

import (
	"context"
	domain "specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/pkg/pagination"
)

// SearchPaymentsFilter represents the filter for searching payments
type SearchPaymentsFilter struct {
	pagination.Paging
}

// PaymentRepository defines the secondary port for payment persistence
type PaymentRepository interface {
	GetAll(ctx context.Context) ([]domain.Payment, error)
	Create(ctx context.Context, payment domain.Payment) (domain.Payment, error)
	SearchPayments(ctx context.Context, filter SearchPaymentsFilter) (pagination.Page[domain.Payment], error)
}
