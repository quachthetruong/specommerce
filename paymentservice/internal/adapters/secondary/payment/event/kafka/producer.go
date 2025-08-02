package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	kafkaGo "github.com/segmentio/kafka-go"
	"specommerce/paymentservice/config"
	domain "specommerce/paymentservice/internal/core/domain/payment"
	"specommerce/paymentservice/internal/core/ports/secondary"
	"specommerce/paymentservice/model"
	"specommerce/paymentservice/pkg/messagequeue"
)

type paymentPublisher struct {
	config    config.AppConfig
	publisher messagequeue.Publisher
}

func (p *paymentPublisher) SendPaymentResponse(ctx context.Context, input domain.ProcessPaymentResponse) error {
	errTemplate := "paymentPublisher SendPaymentResponse failed: %v"

	data := &model.ProcessPaymentResponse{
		PaymentId:     input.PaymentId.String(),
		OrderId:       input.OrderId.String(),
		TotalAmount:   input.TotalAmount,
		CustomerId:    input.CustomerId,
		PaymentStatus: input.Status.String(),
	}
	payload, err := proto.Marshal(data)

	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	return p.publisher.Publish(kafkaGo.Message{
		Topic: p.config.ProcessPaymentResponse.Topic,
		Value: payload,
		Key:   []byte(input.CustomerId),
	})
}

func NewPaymentPublisher(config config.AppConfig, publisher messagequeue.Publisher) secondary.PaymentEventRepository {
	return &paymentPublisher{
		config:    config,
		publisher: publisher,
	}
}
