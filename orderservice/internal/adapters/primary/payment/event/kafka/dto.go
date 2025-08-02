package kafka

import (
	"specommerce/orderservice/internal/core/domain/payment"
)

func ToDomainPaymentStatus(status string) payment.PaymentStatus {
	switch status {
	case "SUCCESS":
		return payment.PaymentStatusSuccess
	case "FAILED":
		return payment.PaymentStatusFailed
	default:
		return payment.PaymentStatusFailed
	}
}
