package payment

import (
	"context"
	"fmt"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/atomicity"
)

type paymentService struct {
	orderRepository secondary.OrderRepository
	atomicExecutor  atomicity.AtomicExecutor
}

func NewPaymentService(
	orderRepository secondary.OrderRepository,
	atomicExecutor atomicity.AtomicExecutor,
) primary.PaymentService {
	return &paymentService{
		orderRepository: orderRepository,
		atomicExecutor:  atomicExecutor,
	}
}

func (s *paymentService) ProcessPaymentResponse(ctx context.Context, input payment.ProcessPaymentResponse) (order.Order, error) {
	errTemplate := "paymentService ProcessPaymentResponse %w"
	orderResponse := order.Order{}
	txErr := s.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			newStatus := order.OrderStatusCompleted
			if input.PaymentStatus == payment.PaymentStatusFailed {
				newStatus = order.OrderStatusFailed
			}
			updatedOrder, err := s.orderRepository.UpdateStatusById(ctx, input.OrderId, newStatus)
			if err != nil {
				return err
			}
			orderResponse = updatedOrder
			return nil
		},
	)
	if txErr != nil {
		return order.Order{}, fmt.Errorf(errTemplate, txErr)
	}
	return orderResponse, nil
}
