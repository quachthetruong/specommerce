package di

import (
	"github.com/samber/do/v2"
	"log/slog"
	"specommerce/paymentservice/config"
	paymentConsumer "specommerce/paymentservice/internal/adapters/primary/payment/event/kafka"
	paymentHandler "specommerce/paymentservice/internal/adapters/primary/payment/handler"
	paymentKafka "specommerce/paymentservice/internal/adapters/secondary/payment/event/kafka"
	paymentPostgres "specommerce/paymentservice/internal/adapters/secondary/payment/persistence/postgres"
	"specommerce/paymentservice/internal/core/ports/primary"
	"specommerce/paymentservice/internal/core/ports/secondary"
	paymentService "specommerce/paymentservice/internal/core/services/payment"
	"specommerce/paymentservice/pkg/atomicity"
	"specommerce/paymentservice/pkg/database"
	"specommerce/paymentservice/pkg/messagequeue"
	"specommerce/paymentservice/pkg/shutdown"
)

func NewInjector() do.Injector {
	injector := do.New()
	do.Provide(injector, NewPaymentRepository)
	do.Provide(injector, NewPaymentService)
	do.Provide(injector, NewPaymentHandler)

	do.Provide(injector, NewPaymentPublisher)
	do.Provide(injector, NewPublisher)
	do.Provide(injector, NewProcessPaymentRequestConsumer)

	do.Provide(injector, NewBaseEventListener)

	return injector
}

func NewPaymentRepository(injector do.Injector) (secondary.PaymentRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](injector)
	return paymentPostgres.NewPaymentPersistenceRepository(getDbFunc), nil
}

func NewPaymentService(injector do.Injector) (primary.PaymentService, error) {
	paymentRepository := do.MustInvoke[secondary.PaymentRepository](injector)
	paymentPublisher := do.MustInvoke[secondary.PaymentEventRepository](injector)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](injector)
	return paymentService.NewPaymentService(
		paymentRepository,
		paymentPublisher,
		atomicExecutor,
	), nil
}

func NewPaymentHandler(injector do.Injector) (paymentHandler.PaymentHandler, error) {
	service := do.MustInvoke[primary.PaymentService](injector)
	return paymentHandler.NewPaymentHandler(service), nil
}

func NewPublisher(injector do.Injector) (messagequeue.Publisher, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	return messagequeue.NewPublisher(cfg, tasks), nil
}

func NewPaymentPublisher(injector do.Injector) (secondary.PaymentEventRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	publisher := do.MustInvoke[messagequeue.Publisher](injector)
	return paymentKafka.NewPaymentPublisher(cfg, publisher), nil
}

func NewBaseEventListener(injector do.Injector) (*messagequeue.BaseEventListener, error) {
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	logger := do.MustInvoke[*slog.Logger](injector)
	return messagequeue.NewBaseEventListener(tasks, logger), nil
}

func NewProcessPaymentRequestConsumer(injector do.Injector) (*paymentConsumer.ProcessPaymentRequestConsumer, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	baseEventListener := do.MustInvoke[*messagequeue.BaseEventListener](injector)
	service := do.MustInvoke[primary.PaymentService](injector)
	return paymentConsumer.NewProcessPaymentRequestConsumer(baseEventListener, cfg.ProcessPaymentRequest, service), nil
}
