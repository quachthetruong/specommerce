package kafka

import (
	"log/slog"

	"github.com/segmentio/kafka-go"
	"specommerce/orderservice/pkg/messagequeue"
	"specommerce/orderservice/pkg/service_config"
)

type ProcessPaymentResponseConsumer struct {
	baseListener *messagequeue.BaseEventListener
	config       service_config.KafkaConfig
}

func NewProcessPaymentResponseConsumer(
	baseListener *messagequeue.BaseEventListener,
	cfg service_config.KafkaConfig,
) *ProcessPaymentResponseConsumer {
	return &ProcessPaymentResponseConsumer{
		baseListener: baseListener,
		config:       cfg,
	}
}

func (c *ProcessPaymentResponseConsumer) Start() error {
	return c.baseListener.Start(c.config, c.handleEvent)
}

func (c *ProcessPaymentResponseConsumer) handleEvent(message kafka.Message) error {
	// TODO: Implement payment response processing logic
	// This should deserialize the payment response message and update order status accordingly

	// Log the received message
	c.baseListener.Logger().Info("Received payment response",
		slog.String("topic", message.Topic),
		slog.String("key", string(message.Key)),
	)

	return nil
}
