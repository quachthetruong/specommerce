package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	kafkaGo "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"specommerce/orderservice/config"
	"specommerce/orderservice/internal/core/domain/order"
	"specommerce/orderservice/internal/core/ports/secondary"
	"specommerce/orderservice/model"
	"specommerce/orderservice/pkg/messagequeue"
)

type campaignPublisher struct {
	config    config.AppConfig
	publisher messagequeue.Publisher
}

func (p *campaignPublisher) SendOrderEvent(ctx context.Context, input order.Order) error {
	errTemplate := "campaignPublisher SendOrderEvent failed: %v"

	payload, err := proto.Marshal(&model.Order{
		Id:           input.Id.String(),
		TotalAmount:  input.TotalAmount,
		CustomerId:   input.CustomerId,
		CustomerName: input.CustomerName,
		Status:       input.Status.String(),
		CreatedAt:    timestamppb.New(input.CreatedAt),
		UpdatedAt:    timestamppb.New(input.UpdatedAt),
	})

	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	return p.publisher.Publish(kafkaGo.Message{
		Topic: p.config.OrderSuccess.Topic,
		Value: payload,
		Key:   []byte(input.CustomerId),
	})
}

func NewCampaignPublisher(config config.AppConfig, publisher messagequeue.Publisher) secondary.CampaignRepository {
	return &campaignPublisher{
		config:    config,
		publisher: publisher,
	}
}
