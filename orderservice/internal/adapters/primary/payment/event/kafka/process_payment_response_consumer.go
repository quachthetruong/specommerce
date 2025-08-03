package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/rs/xid"
	"log/slog"
	"specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/model"

	"github.com/segmentio/kafka-go"
	"specommerce/orderservice/pkg/messagequeue"
	"specommerce/orderservice/pkg/service_config"
)

type ProcessPaymentResponseConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	service      primary.OrderService
}

func NewProcessPaymentResponseConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	service primary.OrderService,
) *ProcessPaymentResponseConsumer {
	return &ProcessPaymentResponseConsumer{
		baseListener: baseListener,
		config:       cfg,
		service:      service,
	}
}

func (c *ProcessPaymentResponseConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *ProcessPaymentResponseConsumer) handleEvent(message kafka.Message) error {
	errorTemplate := "ProcessPaymentResponseConsumer.handleEvent: %w"
	c.baseListener.Logger().Info("Received payment response",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	var request model.ProcessPaymentResponse
	if err := proto.Unmarshal(message.Value, &request); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	ctx := context.Background()
	orderId, err := xid.FromString(request.OrderId)
	if err != nil {
		return fmt.Errorf(errorTemplate, fmt.Errorf(errorTemplate, err))
	}
	successPayment, err := c.service.ProcessPaymentResponse(ctx, payment.ProcessPaymentResponse{
		OrderId:       orderId,
		PaymentStatus: ToDomainPaymentStatus(request.PaymentStatus),
	})

	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	c.baseListener.Logger().Info("Processed payment request successfully",
		slog.String("payment_id", request.PaymentId),
		slog.String("order_id", successPayment.Id.String()),
		slog.Float64("total_amount", successPayment.TotalAmount),
		slog.String("customer_id", successPayment.CustomerId),
		slog.String("status", successPayment.Status.String()),
	)

	return nil
}
