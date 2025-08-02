package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/model"

	"github.com/segmentio/kafka-go"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/service_config"
)

type OrderCompletedConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	orderService primary.OrderService
}

func NewOrderCompletedConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	orderService primary.OrderService,
) *OrderCompletedConsumer {
	return &OrderCompletedConsumer{
		baseListener: baseListener,
		config:       cfg,
		orderService: orderService,
	}
}

func (c *OrderCompletedConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *OrderCompletedConsumer) handleEvent(message kafka.Message) error {
	errorTemplate := "OrderCompletedConsumer.handleEvent: %w"
	c.baseListener.Logger().Info("Received order complete event",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	// Parse the order complete event
	var orderEvent model.OrderCompleted
	if err := json.Unmarshal(message.Value, &orderEvent); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	// Convert to domain object
	domainOrder, err := ToDomain(orderEvent)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	ctx := context.Background()

	// Process the completed order for campaign tracking
	savedOrder, err := c.orderService.CreateOrder(ctx, domainOrder)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	c.baseListener.Logger().Info("Processed order completed event successfully",
		slog.String("order_id", savedOrder.Id.String()),
		slog.String("customer_id", domainOrder.CustomerId),
		slog.Float64("total_amount", domainOrder.TotalAmount),
		slog.String("status", domainOrder.Status.String()),
	)

	return nil
}
