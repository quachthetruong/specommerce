package handler

import (
	domain "specommerce/paymentservice/internal/core/domain/payment"
	"time"
)

func ToGetAllPaymentResponse(entities []domain.Payment) []PaymentResponse {
	response := make([]PaymentResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, PaymentResponse{
			ID:          entity.Id.String(),
			OrderID:     entity.OrderId.String(),
			CustomerID:  entity.CustomerId,
			TotalAmount: entity.TotalAmount,
			Status:      entity.Status.String(),
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		})
	}
	return response
}

// PaymentResponse represents payment response for Swagger
type PaymentResponse struct {
	ID          string    `json:"id" example:"abc123"`
	OrderID     string    `json:"order_id" example:"order123"`
	CustomerID  string    `json:"customer_id" example:"customer123"`
	TotalAmount float64   `json:"total_amount" example:"99.99"`
	Status      string    `json:"status" example:"SUCCESS"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}
