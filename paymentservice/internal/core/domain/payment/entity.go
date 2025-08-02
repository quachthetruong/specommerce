package payment

import (
	"github.com/rs/xid"
	"time"
)

type ProcessPaymentResponse struct {
	PaymentId   xid.ID        `json:"payment_id" validate:"required"`
	OrderId     xid.ID        `json:"order_id" validate:"required"`
	Status      PaymentStatus `json:"status" validate:"required"`
	CustomerId  string        `json:"customer_id" validate:"required"`
	TotalAmount float64       `json:"total_amount" validate:"required,gt=0"`
}

type ProcessPaymentRequest struct {
	OrderId     xid.ID  `json:"order_id" validate:"required"`
	CustomerId  string  `json:"customer_id" validate:"required"`
	TotalAmount float64 `json:"total_amount" validate:"required,gt=0"`
}
type PaymentStatus string

const (
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	Id          xid.ID        `json:"id" bun:"id,pk,skipupdate"`
	OrderId     xid.ID        `json:"order_id" bun:"order_id,notnull"` // Reference to the order
	CustomerId  string        `json:"customer_id" bun:"customer_id"`
	TotalAmount float64       `json:"total_amount" bun:"total_amount"`
	Status      PaymentStatus `json:"status" bun:"status"`
	CreatedAt   time.Time     `json:"created_at" bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt   time.Time     `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}

func (s PaymentStatus) String() string {
	return string(s)
}
