package handler

import (
	"github.com/rs/xid"
	domain "specommerce/orderservice/internal/core/domain/order"
	"time"
)

// CreateOrderRequest represents the request for creating an order
type CreateOrderRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
}

// ToOrder converts CreateOrderRequest to domain Order
func (r *CreateOrderRequest) ToDomain() domain.Order {
	return domain.Order{
		Id:          xid.New(),
		CustomerId:  r.CustomerID,
		Status:      domain.OrderStatusPending,
		TotalAmount: r.TotalAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func ToCreateOrderResponse(d domain.Order) OrderResponse {
	return OrderResponse{
		ID:          d.Id.String(),
		CustomerID:  d.CustomerId,
		Status:      d.Status.String(),
		TotalAmount: d.TotalAmount,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func ToGetAllOrderResponse(entities []domain.Order) []OrderResponse {
	response := make([]OrderResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, OrderResponse{
			ID:          entity.Id.String(),
			CustomerID:  entity.CustomerId,
			Status:      entity.Status.String(),
			TotalAmount: entity.TotalAmount,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		})
	}
	return response

}

// OrderResponse represents order response for Swagger
type OrderResponse struct {
	ID          string    `json:"id" example:"abc123"`
	CustomerID  string    `json:"customer_id" example:"customer123"`
	TotalAmount float64   `json:"total_amount" example:"99.99"`
	Status      string    `json:"status" example:"PENDING"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}
