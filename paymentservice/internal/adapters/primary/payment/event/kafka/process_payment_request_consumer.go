package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/rs/xid"
	"log/slog"
	"specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/internal/core/ports/primary"
	"specommerce/paymentservice/model"
	"time"

	"github.com/segmentio/kafka-go"
	"specommerce/paymentservice/pkg/messagequeue"
	"specommerce/paymentservice/pkg/service_config"
)

type ProcessPaymentRequestConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
	service      primary.PaymentService
}

func NewProcessPaymentRequestConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
	service primary.PaymentService,
) *ProcessPaymentRequestConsumer {
	return &ProcessPaymentRequestConsumer{
		baseListener: baseListener,
		config:       cfg,
		service:      service,
	}
}

func (c *ProcessPaymentRequestConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *ProcessPaymentRequestConsumer) handleEvent(message kafka.Message) error {
	errorTemplate := "ProcessPaymentRequestConsumer.handleEvent: %w"
	c.baseListener.Logger().Info("Received payment response",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	var request model.ProcessPaymentRequest
	if err := proto.Unmarshal(message.Value, &request); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	ctx := context.Background()
	time.Sleep(time.Duration(request.TimeProcess) * time.Millisecond)
	orderId, err := xid.FromString(request.OrderId)
	if err != nil {
		return fmt.Errorf(errorTemplate, fmt.Errorf("invalid order ID: %w", err))
	}
	successPayment, err := c.service.ProcessPaymentRequest(ctx, payment.Payment{
		Id:          xid.New(),
		OrderId:     orderId,
		TotalAmount: request.TotalAmount,
		CustomerId:  request.CustomerId,
		Status:      payment.PaymentStatusSuccess,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	c.baseListener.Logger().Info("Processed payment request successfully",
		slog.String("payment_id", successPayment.Id.String()),
		slog.String("order_id", successPayment.OrderId.String()),
		slog.Float64("total_amount", successPayment.TotalAmount),
		slog.String("customer_id", successPayment.CustomerId),
		slog.String("status", successPayment.Status.String()),
	)

	return nil
}
