package postgres

import (
	"github.com/rs/xid"
	"github.com/uptrace/bun"
	domain "specommerce/paymentservice/internal/core/domain/payment"
	"time"
)

type Payment struct {
	bun.BaseModel `bun:"payments"`
	Id            xid.ID    `bun:",skipupdate,pk"`
	OrderId       xid.ID    `bun:"order_id,notnull"` // Reference to the order
	TotalAmount   float64   `bun:"total_amount,notnull"`
	CustomerId    string    `bun:"customer_id,notnull"`
	Status        string    `bun:"status,notnull,default:'SUCCESS'"` //
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (o Payment) ToDomainModel() domain.Payment {
	return domain.Payment{
		Id:          o.Id,
		OrderId:     o.OrderId,
		CustomerId:  o.CustomerId,
		TotalAmount: o.TotalAmount,
		Status:      domain.PaymentStatus(o.Status),
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func FromDomainModel(dm domain.Payment) Payment {
	return Payment{
		Id:          dm.Id,
		OrderId:     dm.OrderId,
		CustomerId:  dm.CustomerId,
		TotalAmount: dm.TotalAmount,
		Status:      string(dm.Status),
		CreatedAt:   dm.CreatedAt,
		UpdatedAt:   dm.UpdatedAt,
	}
}
