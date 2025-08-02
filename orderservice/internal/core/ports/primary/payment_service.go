package primary

import (
	"context"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
)

type PaymentService interface {
	ProcessPaymentResponse(ctx context.Context, request payment.ProcessPaymentResponse) (order.Order, error)
}
