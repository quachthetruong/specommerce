package order

import (
	"context"
	"fmt"
	"specommerce/campaignservice/internal/core/domain/customer"
	"specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	"specommerce/campaignservice/pkg/atomicity"
)

// OrderService implements the order business logic
type service struct {
	orderRepo      secondary.OrderRepository
	atomicExecutor atomicity.AtomicExecutor
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo secondary.OrderRepository, atomicExecutor atomicity.AtomicExecutor) primary.OrderService {
	return &service{
		orderRepo:      orderRepo,
		atomicExecutor: atomicExecutor,
	}
}

// CreateOrder creates a new order
func (s *service) CreateOrder(ctx context.Context, input order.Order) (order.Order, error) {
	errTemplate := "orderService CreateOrder %w"
	savedOrder, err := s.orderRepo.Create(ctx, input)
	if err != nil {
		return order.Order{}, fmt.Errorf(errTemplate, err)
	}
	return savedOrder, nil
}

func (s *service) GetWinner(ctx context.Context, campaignID int64) ([]customer.Customer, error) {
	return []customer.Customer{}, nil
}
