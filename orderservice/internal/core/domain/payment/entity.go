package payment

import "github.com/rs/xid"

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type ProcessPaymentRequest struct {
	OrderId     xid.ID  `json:"order_id" validate:"required"`
	CustomerId  string  `json:"customer_id" validate:"required"`
	TotalAmount float64 `json:"total_amount" validate:"required,gt=0"`
}

type ProcessPaymentResponse struct {
	OrderId       xid.ID        `json:"order_id" validate:"required"`
	PaymentStatus PaymentStatus `json:"payment_status" validate:"required"`
}
