package payment

type ProcessPaymentRequest struct {
	OrderId     string  `json:"order_id" validate:"required"`
	CustomerId  string  `json:"customer_id" validate:"required"`
	TotalAmount float64 `json:"total_amount" validate:"required,gt=0"`
}
