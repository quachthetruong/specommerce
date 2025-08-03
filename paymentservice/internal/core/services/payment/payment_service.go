package payment

import (
	"context"
	"fmt"
	"specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/internal/core/ports/primary"
	"specommerce/paymentservice/internal/core/ports/secondary"
	"specommerce/paymentservice/pkg/atomicity"
	"specommerce/paymentservice/pkg/pagination"
)

type paymentService struct {
	paymentRepository secondary.PaymentRepository
	paymentPublisher  secondary.PaymentEventRepository
	atomicExecutor    atomicity.AtomicExecutor
}

func NewPaymentService(
	paymentRepository secondary.PaymentRepository,
	paymentPublisher secondary.PaymentEventRepository,
	atomicExecutor atomicity.AtomicExecutor,
) primary.PaymentService {
	return &paymentService{
		paymentRepository: paymentRepository,
		paymentPublisher:  paymentPublisher,
		atomicExecutor:    atomicExecutor,
	}
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]payment.Payment, error) {
	return s.paymentRepository.GetAll(ctx)
}

func (s *paymentService) ProcessPaymentRequest(ctx context.Context, input payment.Payment) (payment.Payment, error) {
	errTemplate := "paymentService ProcessPaymentRequest %w"
	paymentResponse := payment.Payment{}
	txErr := s.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			pendingOrder, err := s.paymentRepository.Create(ctx, input)
			if err != nil {
				return err
			}

			err = s.paymentPublisher.SendPaymentResponse(ctx, payment.ProcessPaymentResponse{
				PaymentId:   pendingOrder.Id,
				OrderId:     pendingOrder.OrderId,
				CustomerId:  pendingOrder.CustomerId,
				TotalAmount: pendingOrder.TotalAmount,
				Status:      pendingOrder.Status,
			})
			if err != nil {
				return err
			}
			paymentResponse = pendingOrder
			return nil
		},
	)
	if txErr != nil {
		return payment.Payment{}, fmt.Errorf(errTemplate, txErr)
	}
	return paymentResponse, nil

}

func (s *paymentService) SearchPayments(ctx context.Context, filter secondary.SearchPaymentsFilter) (pagination.Page[payment.Payment], error) {
	return s.paymentRepository.SearchPayments(ctx, filter)
}
