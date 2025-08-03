package primary

import (
	"context"
	"specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/internal/core/ports/secondary"
	"specommerce/paymentservice/pkg/pagination"
)

// PaymentService defines the primary port for payment operations
type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]payment.Payment, error)
	ProcessPaymentRequest(ctx context.Context, input payment.Payment) (payment.Payment, error)
	SearchPayments(ctx context.Context, filter secondary.SearchPaymentsFilter) (pagination.Page[payment.Payment], error)
}
