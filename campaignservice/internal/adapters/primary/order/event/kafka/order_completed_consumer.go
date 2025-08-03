package kafka

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/model"

	"github.com/segmentio/kafka-go"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/service_config"
)

type OrderSuccessConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	orderService primary.OrderService
}

func NewOrderSuccessConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	orderService primary.OrderService,
) *OrderSuccessConsumer {
	return &OrderSuccessConsumer{
		baseListener: baseListener,
		config:       cfg,
		orderService: orderService,
	}
}

func (c *OrderSuccessConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *OrderSuccessConsumer) handleEvent(message kafka.Message) error {
	errorTemplate := "OrderSuccessConsumer.handleEvent: %w"
	c.baseListener.Logger().Info("Received order success event",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	var orderEvent model.Order
	if err := proto.Unmarshal(message.Value, &orderEvent); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	domainOrder, err := ToDomain(orderEvent)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	ctx := context.Background()

	savedOrder, err := c.orderService.CreateOrder(ctx, domainOrder)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	c.baseListener.Logger().Info("Processed order success event successfully",
		slog.String("order_id", savedOrder.Id.String()),
		slog.String("customer_id", domainOrder.CustomerId),
		slog.Float64("total_amount", domainOrder.TotalAmount),
		slog.String("status", domainOrder.Status.String()),
	)

	return nil
}
