package order

import (
	"time"

	"github.com/rs/xid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
)

type Order struct {
	ID          xid.ID      `json:"id" bun:"id,pk,skipupdate"`
	CustomerID  string      `json:"customer_id" bun:"customer_id"`
	TotalAmount float64     `json:"total_amount" bun:"total_amount"`
	Status      OrderStatus `json:"status" bun:"status"`
	CreatedAt   time.Time   `json:"created_at" bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt   time.Time   `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}

func (s OrderStatus) String() string {
	return string(s)
}
