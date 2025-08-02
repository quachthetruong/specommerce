package secondary

import (
	"context"
	domain "specommerce/paymentservice/internal/core/domain/payment"
)

// PaymentRepository defines the secondary port for payment persistence
type PaymentRepository interface {
	GetAll(ctx context.Context) ([]domain.Payment, error)
	Create(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
