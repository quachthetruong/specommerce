package messagequeue

import (
	"context"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
	"specommerce/paymentservice/pkg/service_config"
	"specommerce/paymentservice/pkg/shutdown"
)

type EventListener interface {
	Start() error
}

type HandlerFunc func(message kafka.Message) error

type BaseEventListener struct {
	logger       *slog.Logger
	shutdownTask *shutdown.Tasks
}

func NewBaseEventListener(shutdownTask *shutdown.Tasks, logger *slog.Logger) *BaseEventListener {
	return &BaseEventListener{
		shutdownTask: shutdownTask,
		logger:       logger,
	}
}

func (l *BaseEventListener) Logger() *slog.Logger {
	return l.logger
}

func (l *BaseEventListener) Start(cfg service_config.KafkaConfig, handlerFunc HandlerFunc) error {
	l.logger.Info("Starting Kafka event listener",
		slog.String("brokers", cfg.Host),
		slog.String("consumer_group", cfg.ConsumerGroup),
		slog.Any("topic", cfg.Topic),
	)

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{cfg.Host},
			Topic:   cfg.Topic,
			GroupID: cfg.ConsumerGroup,
		},
	)
	loop := true
	l.shutdownTask.AddShutdownTask(
		func(ctx context.Context) error {
			loop = false
			return reader.Close()
		},
	)
	var waitGroup sync.WaitGroup
	l.shutdownTask.AddShutdownTask(
		func(ctx context.Context) error {
			waitGroup.Wait()
			return nil
		},
	)
	for loop {
		//NOTE: this auto committed after fetch message -> need to re-consider
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			l.logger.Error("Failed to read message from Kafka", slog.String("error", err.Error()))
			continue
		}
		waitGroup.Add(1)
		go func(msg kafka.Message) {
			defer waitGroup.Done()
			err := handlerFunc(msg)
			if err != nil {
				l.logger.Error("Can not handle event",
					slog.String("error", err.Error()),
					slog.String("topic", msg.Topic),
					slog.String("key", string(msg.Key)),
				)
			}
		}(message)
	}
	return nil
}
