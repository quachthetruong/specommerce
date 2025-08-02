package primary

import (
	"context"
	"specommerce/paymentservice/internal/core/domain/payment"
)

// PaymentService defines the primary port for payment operations
type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]payment.Payment, error)
	ProcessPaymentRequest(ctx context.Context, input payment.Payment) (payment.Payment, error)
}
