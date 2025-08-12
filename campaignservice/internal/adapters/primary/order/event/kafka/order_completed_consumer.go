package kafka

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log/slog"
	domain "specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/model"

	"github.com/segmentio/kafka-go"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/service_config"
)

type OrderConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	orderService primary.OrderService
}

func NewOrderConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	orderService primary.OrderService,
) *OrderConsumer {
	return &OrderConsumer{
		baseListener: baseListener,
		config:       cfg,
		orderService: orderService,
	}
}

func (c *OrderConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *OrderConsumer) handleEvent(message kafka.Message) error {
	errorTemplate := "OrderConsumer.handleEvent: %w"
	c.baseListener.Logger().Info("Received order event",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	var orderEvent model.Order
	if err := proto.Unmarshal(message.Value, &orderEvent); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	order, err := ToDomain(orderEvent)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	ctx := context.Background()

	if order.Status == domain.OrderStatusPending {
		err = c.orderService.ProcessPendingOrder(ctx, order)
		if err != nil {
			return fmt.Errorf(errorTemplate, err)
		}
		c.baseListener.Logger().Info("Processed pending order",
			slog.String("order_id", order.Id.String()),
			slog.String("customer_id", order.CustomerId),
			slog.Float64("total_amount", order.TotalAmount),
			slog.String("status", order.Status.String()),
		)
		return nil
	}
	err = c.orderService.ProcessOrderResult(ctx, order)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	c.baseListener.Logger().Info("Processed order event successfully",
		slog.String("order_id", order.Id.String()),
		slog.String("customer_id", order.CustomerId),
		slog.Float64("total_amount", order.TotalAmount),
		slog.String("status", order.Status.String()),
	)

	return nil
}
