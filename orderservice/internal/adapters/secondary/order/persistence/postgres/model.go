package postgres

import (
	"github.com/rs/xid"
	"github.com/uptrace/bun"
	domain "specommerce/orderservice/internal/core/domain/order"
	"time"
)

type Order struct {
	bun.BaseModel `bun:"orders"`
	Id            xid.ID    `bun:",skipupdate,pk"`
	TotalAmount   float64   `bun:"total_amount,notnull"`
	CustomerId    string    `bun:"customer_id,notnull"`
	CustomerName  string    `bun:"customer_name,notnull"`            // Added field for customer name
	Status        string    `bun:"status,notnull,default:'PENDING'"` //
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (o Order) ToDomainModel() domain.Order {
	return domain.Order{
		Id:           o.Id,
		CustomerId:   o.CustomerId,
		CustomerName: o.CustomerName,
		TotalAmount:  o.TotalAmount,
		Status:       domain.OrderStatus(o.Status),
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

func FromDomainModel(dm domain.Order) Order {
	return Order{
		Id:           dm.Id,
		CustomerId:   dm.CustomerId,
		CustomerName: dm.CustomerName,
		TotalAmount:  dm.TotalAmount,
		Status:       string(dm.Status),
		CreatedAt:    dm.CreatedAt,
		UpdatedAt:    dm.UpdatedAt,
	}
}
