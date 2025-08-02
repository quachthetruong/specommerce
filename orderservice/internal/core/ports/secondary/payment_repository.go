package secondary

import (
	"context"
	domain "specommerce/orderservice/internal/core/domain/payment"
)

type PaymentRepository interface {
	SendPaymentRequest(ctx context.Context, input domain.ProcessPaymentRequest) error
}
