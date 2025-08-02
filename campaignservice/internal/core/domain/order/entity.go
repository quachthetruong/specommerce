package order

import (
	"time"

	"github.com/rs/xid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusSuccess    OrderStatus = "SUCCESS"
	OrderStatusFailed     OrderStatus = "FAILED"
)

type Order struct {
	Id           xid.ID      `json:"id" bun:"id,pk,skipupdate"`
	CustomerId   string      `json:"customer_id" bun:"customer_id"`
	CustomerName string      `json:"customer_name" bun:"customer_name"`
	TotalAmount  float64     `json:"total_amount" bun:"total_amount"`
	Status       OrderStatus `json:"status" bun:"status"`
	CreatedAt    time.Time   `json:"created_at" bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt    time.Time   `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}

func (s OrderStatus) String() string {
	return string(s)
}
