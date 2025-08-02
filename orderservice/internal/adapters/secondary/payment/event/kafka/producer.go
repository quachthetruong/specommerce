package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	kafkaGo "github.com/segmentio/kafka-go"
	"specommerce/orderservice/config"
	domain "specommerce/orderservice/internal/core/domain/payment"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/model"
	"specommerce/orderservice/pkg/messagequeue"
)

type paymentPublisher struct {
	config    config.AppConfig
	publisher messagequeue.Publisher
}

func (p *paymentPublisher) SendPaymentRequest(ctx context.Context, input domain.ProcessPaymentRequest) error {
	errTemplate := "paymentPublisher PublishProcessPaymentRequest failed: %v"

	sellConfirmationMessage := &model.ProcessPaymentRequest{
		OrderId:     input.OrderId.String(),
		TotalAmount: input.TotalAmount,
		CustomerId:  input.CustomerId,
	}
	payload, err := proto.Marshal(sellConfirmationMessage)

	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	return p.publisher.Publish(kafkaGo.Message{
		Topic: p.config.ProcessPaymentRequest.Topic,
		Value: payload,
		Key:   []byte(input.CustomerId),
	})
}

func NewPaymentPublisher(config config.AppConfig, publisher messagequeue.Publisher) secondary.PaymentRepository {
	return &paymentPublisher{
		config:    config,
		publisher: publisher,
	}
}
