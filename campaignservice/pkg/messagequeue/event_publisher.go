package messagequeue

import (
	"context"

	"github.com/segmentio/kafka-go"

	"specommerce/campaignservice/config"
	"specommerce/campaignservice/pkg/shutdown"
)

type Publisher interface {
	Publish(message kafka.Message) error
}

type publisher struct {
	kafkaWriter *kafka.Writer
}

func NewPublisher(cfg config.AppConfig, tasks *shutdown.Tasks) Publisher {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Host),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           1,
		AllowAutoTopicCreation: cfg.Kafka.AutoCreateTopic,
		MaxAttempts:            cfg.Kafka.Retry,
	}
	tasks.AddShutdownTask(
		func(ctx context.Context) error {
			return writer.Close()
		},
	)
	return &publisher{kafkaWriter: writer}
}

func (publisher *publisher) Publish(message kafka.Message) error {
	return publisher.kafkaWriter.WriteMessages(context.Background(), message)
}
