package di

import (
	"github.com/samber/do/v2"
	"log/slog"
	"specommerce/orderservice/config"
	orderHandler "specommerce/orderservice/internal/adapters/primary/order/handler"
	paymentConsumer "specommerce/orderservice/internal/adapters/primary/payment/event/kafka"
	orderPostgres "specommerce/orderservice/internal/adapters/secondary/order/persistence/postgres"
	paymentKafka "specommerce/orderservice/internal/adapters/secondary/payment/event/kafka"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/internal/core/ports/secondary"
	orderService "specommerce/orderservice/internal/core/services/order"
	"specommerce/orderservice/pkg/atomicity"
	"specommerce/orderservice/pkg/database"
	"specommerce/orderservice/pkg/messagequeue"
	"specommerce/orderservice/pkg/shutdown"
)

func NewInjector() do.Injector {
	injector := do.New()
	do.Provide(injector, NewOrderRepository)
	do.Provide(injector, NewOrderService)
	do.Provide(injector, NewOrderHandler)

	do.Provide(injector, NewPaymentPublisher)
	do.Provide(injector, NewPublisher)
	do.Provide(injector, NewProcessPaymentResponseConsumer)

	do.Provide(injector, NewBaseEventListener)

	return injector
}

func NewOrderRepository(injector do.Injector) (secondary.OrderRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](injector)
	return orderPostgres.NewOrderPersistenceRepository(getDbFunc), nil
}

func NewOrderService(injector do.Injector) (primary.OrderService, error) {
	orderRepository := do.MustInvoke[secondary.OrderRepository](injector)
	paymentPublisher := do.MustInvoke[secondary.PaymentRepository](injector)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](injector)
	return orderService.NewOrderService(
		orderRepository,
		paymentPublisher,
		atomicExecutor,
	), nil
}

func NewOrderHandler(injector do.Injector) (orderHandler.OrderHandler, error) {
	service := do.MustInvoke[primary.OrderService](injector)
	return orderHandler.NewOrderHandler(service), nil
}

func NewPublisher(injector do.Injector) (messagequeue.Publisher, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	return messagequeue.NewPublisher(cfg, tasks), nil
}

func NewPaymentPublisher(injector do.Injector) (secondary.PaymentRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	publisher := do.MustInvoke[messagequeue.Publisher](injector)
	return paymentKafka.NewPaymentPublisher(cfg, publisher), nil
}

func NewBaseEventListener(injector do.Injector) (*messagequeue.BaseEventListener, error) {
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	logger := do.MustInvoke[*slog.Logger](injector)
	return messagequeue.NewBaseEventListener(tasks, logger), nil
}

func NewProcessPaymentResponseConsumer(injector do.Injector) (*paymentConsumer.ProcessPaymentResponseConsumer, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	baseEventListener := do.MustInvoke[*messagequeue.BaseEventListener](injector)
	return paymentConsumer.NewProcessPaymentResponseConsumer(baseEventListener, cfg.ProcessPaymentResponse), nil
}
