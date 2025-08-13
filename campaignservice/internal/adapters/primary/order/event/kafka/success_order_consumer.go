package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"log/slog"
	domain "specommerce/campaignservice/internal/core/domain/order"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/model"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/service_config"
)

type SuccessOrderConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	orderService primary.OrderService
}

func NewSuccessOrderConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	orderService primary.OrderService,
) *SuccessOrderConsumer {
	return &SuccessOrderConsumer{
		baseListener: baseListener,
		config:       cfg,
		orderService: orderService,
	}
}

func (c *SuccessOrderConsumer) Start() error {
	return c.baseListener.Start(c.config, c.HandleEvent)
}

func (c *SuccessOrderConsumer) HandleEvent(message kafka.Message) error {
	errorTemplate := "SuccessOrderConsumer.HandleEvent: %w"
	c.baseListener.Logger().Info("Received success order event",
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

	if order.Status != domain.OrderStatusSuccess {
		return nil
	}
	err = c.orderService.SaveSuccessOrder(ctx, order)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	c.baseListener.Logger().Info("Save success order",
		slog.String("order_id", order.Id.String()),
		slog.String("customer_id", order.CustomerId),
		slog.Float64("total_amount", order.TotalAmount),
		slog.String("status", order.Status.String()),
	)
	return nil
}
