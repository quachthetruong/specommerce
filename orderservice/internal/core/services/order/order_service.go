package order

import (
	"context"
	"fmt"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/atomicity"
)

// OrderService implements the order business logic
type service struct {
	orderRepo        secondary.OrderRepository
	paymentPublisher secondary.PaymentRepository
	atomicExecutor   atomicity.AtomicExecutor
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo secondary.OrderRepository, paymentPublisher secondary.PaymentRepository, atomicExecutor atomicity.AtomicExecutor) primary.OrderService {
	return &service{
		orderRepo:        orderRepo,
		paymentPublisher: paymentPublisher,
		atomicExecutor:   atomicExecutor,
	}
}

// CreateOrder creates a new order
func (s *service) CreateOrder(ctx context.Context, input order.Order) (order.Order, error) {
	errTemplate := "orderService CreateOrder %w"
	orderResponse := order.Order{}
	txErr := s.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			pendingOrder, err := s.orderRepo.Create(ctx, input)
			if err != nil {
				return err
			}

			err = s.paymentPublisher.SendPaymentRequest(ctx, payment.ProcessPaymentRequest{
				OrderId:     pendingOrder.Id,
				CustomerId:  pendingOrder.CustomerId,
				TotalAmount: pendingOrder.TotalAmount,
			})
			if err != nil {
				return err
			}
			orderResponse = pendingOrder
			return nil
		},
	)
	if txErr != nil {
		return order.Order{}, fmt.Errorf(errTemplate, txErr)
	}
	return orderResponse, nil

}

// GetAllOrders retrieves all orders
func (s *service) GetAllOrders(ctx context.Context) ([]order.Order, error) {
	return s.orderRepo.GetAll(ctx)
}
