package kafka

import (
	"fmt"
	"github.com/rs/xid"
	domain "specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/model"
)

func ToDomain(event model.Order) (domain.Order, error) {
	errTemplate := "OrderConsumer.ToDomain: %w"
	orderId, err := xid.FromString(event.Id)
	if err != nil {
		return domain.Order{}, fmt.Errorf(errTemplate, err)
	}
	return domain.Order{
		Id:           orderId,
		CustomerId:   event.CustomerId,
		CustomerName: event.CustomerName,
		TotalAmount:  event.TotalAmount,
		Status:       domain.OrderStatus(event.Status),
		CreatedAt:    event.CreatedAt.AsTime(),
		UpdatedAt:    event.UpdatedAt.AsTime(),
	}, nil

}
