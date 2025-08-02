package kafka

import (
	"specommerce/orderservice/internal/core/domain/payment"
)

func ToDomainPaymentStatus(status string) payment.PaymentStatus {
	switch status {
	case "COMPLETED":
		return payment.PaymentStatusCompleted
	case "FAILED":
		return payment.PaymentStatusFailed
	default:
		return payment.PaymentStatusFailed
	}
}
