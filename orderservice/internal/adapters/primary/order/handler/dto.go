package handler

import (
	"github.com/rs/xid"
	domain "specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/pagination"
	"time"
)

// CreateOrderRequest represents the request for creating an order
type CreateOrderRequest struct {
	CustomerId   string  `json:"customer_id" binding:"required"`
	CustomerName string  `json:"customer_name" binding:"required"`
	TotalAmount  float64 `json:"total_amount" binding:"required"`
	TimeProcess  int64   `json:"time_process" binding:"required" default:"2"`
}

// ToOrder converts CreateOrderRequest to domain Order
func (r *CreateOrderRequest) ToDomain() domain.CreateOrderRequest {
	return domain.CreateOrderRequest{
		Order: domain.Order{
			Id:           xid.New(),
			CustomerId:   r.CustomerId,
			CustomerName: r.CustomerName,
			Status:       domain.OrderStatusPending,
			TotalAmount:  r.TotalAmount,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		TimeProcess: r.TimeProcess,
	}
}

func ToCreateOrderResponse(d domain.Order) OrderResponse {
	return OrderResponse{
		ID:           d.Id.String(),
		CustomerId:   d.CustomerId,
		CustomerName: d.CustomerName,
		Status:       d.Status.String(),
		TotalAmount:  d.TotalAmount,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func ToGetAllOrderResponse(entities []domain.Order) []OrderResponse {
	response := make([]OrderResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, OrderResponse{
			ID:           entity.Id.String(),
			CustomerId:   entity.CustomerId,
			CustomerName: entity.CustomerName,
			Status:       entity.Status.String(),
			TotalAmount:  entity.TotalAmount,
			CreatedAt:    entity.CreatedAt,
			UpdatedAt:    entity.UpdatedAt,
		})
	}
	return response

}

// OrderResponse represents order response for Swagger
type OrderResponse struct {
	ID           string    `json:"id" example:"abc123"`
	CustomerId   string    `json:"customer_id" example:"customer123"`
	CustomerName string    `json:"customer_name" example:"John Doe"`
	TotalAmount  float64   `json:"total_amount" example:"99.99"`
	Status       string    `json:"status" example:"PENDING"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SearchOrdersRequest represents the request for searching orders with pagination
type SearchOrdersRequest struct {
	Paging pagination.Paging
}

func (req SearchOrdersRequest) ToFilter() secondary.SearchOrdersFilter {
	return secondary.SearchOrdersFilter{
		Paging: req.Paging,
	}
}
