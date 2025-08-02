package secondary

import (
	"context"
	domain "specommerce/paymentservice/internal/core/domain/payment"
)

type PaymentEventRepository interface {
	SendPaymentResponse(ctx context.Context, input domain.ProcessPaymentResponse) error
}
