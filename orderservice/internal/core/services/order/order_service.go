package order

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"log/slog"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/pkg/atomicity"
)

// OrderService implements the order business logic
type service struct {
	orderRepo         secondary.OrderRepository
	paymentPublisher  secondary.PaymentRepository
	campaignPublisher secondary.CampaignRepository
	atomicExecutor    atomicity.AtomicExecutor
	logger            *slog.Logger
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo secondary.OrderRepository, paymentPublisher secondary.PaymentRepository, atomicExecutor atomicity.AtomicExecutor,
	campaignPublisher secondary.CampaignRepository, logger *slog.Logger) primary.OrderService {
	return &service{
		orderRepo:         orderRepo,
		campaignPublisher: campaignPublisher,
		paymentPublisher:  paymentPublisher,
		atomicExecutor:    atomicExecutor,
		logger:            logger,
	}
}

// CreateOrder creates a new order and initiates payment processing
// Step 1: Create the order in the database with status Pending
// Step 2: Send a payment request to the payment service
// Step 3: If sending the payment request is successful, update the order status to Processing
// TODO: Put all the steps in a workflow or saga pattern
func (s *service) CreateOrder(ctx context.Context, input order.Order) (order.Order, error) {
	errTemplate := "orderService CreateOrder %w"
	var orderId xid.ID
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
			orderId = pendingOrder.Id
			return nil
		},
	)
	if txErr != nil {
		return order.Order{}, fmt.Errorf(errTemplate, txErr)
	}
	processingOrder, err := s.orderRepo.UpdateStatusById(ctx, orderId, order.OrderStatusProcessing)
	if err != nil {
		return order.Order{}, fmt.Errorf(errTemplate, err)
	}
	return processingOrder, nil

}

// ProcessPaymentResponse processes the payment response from the payment service
// Step 1: Check the payment status
// Step 2: If the payment is successful, update the order status to Success. If the payment failed, update the order status to Failed
// Step 3: If the order status is Success, send an order success event to the campaign service
// TODO: Use CDC to decouple the campaign service logic from the order service
func (s *service) ProcessPaymentResponse(ctx context.Context, input payment.ProcessPaymentResponse) (order.Order, error) {
	errTemplate := "paymentService ProcessPaymentResponse %w"
	orderResponse := order.Order{}
	txErr := s.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			newStatus := order.OrderStatusSuccess
			if input.PaymentStatus == payment.PaymentStatusFailed {
				newStatus = order.OrderStatusFailed
			}
			updatedOrder, err := s.orderRepo.UpdateStatusById(ctx, input.OrderId, newStatus)
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
	if orderResponse.Status != order.OrderStatusSuccess {
		return orderResponse, nil
	}
	go func() {
		err := s.campaignPublisher.SendOrderSuccessEvent(ctx, orderResponse)
		if err != nil {
			s.logger.Error(
				"failed to send order success event to campaign service",
				slog.String("order_id", orderResponse.Id.String()),
				slog.String("error", err.Error()),
			)
		}
	}()

	return orderResponse, nil
}

func (s *service) GetAllOrders(ctx context.Context) ([]order.Order, error) {
	return s.orderRepo.GetAll(ctx)
}
