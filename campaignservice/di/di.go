package di

import (
	"github.com/samber/do/v2"
	"log/slog"
	"specommerce/campaignservice/config"
	campaignHandler "specommerce/campaignservice/internal/adapters/primary/campaign/handler"
	orderConsumer "specommerce/campaignservice/internal/adapters/primary/order/event/kafka"
	campaignPostgres "specommerce/campaignservice/internal/adapters/secondary/campaign/persistence/postgres"
	orderPostgres "specommerce/campaignservice/internal/adapters/secondary/order/persistence/postgres"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/internal/core/ports/secondary"
	campaignService "specommerce/campaignservice/internal/core/services/campaign"
	orderService "specommerce/campaignservice/internal/core/services/order"

	"specommerce/campaignservice/pkg/atomicity"
	"specommerce/campaignservice/pkg/database"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/shutdown"
)

func NewInjector() do.Injector {
	injector := do.New()

	// Core dependencies
	do.Provide(injector, NewCampaignRepository)
	do.Provide(injector, NewCampaignService)
	do.Provide(injector, NewCampaignHandler)

	do.Provide(injector, NewOrderRepository)
	do.Provide(injector, NewOrderService)

	// Messaging dependencies
	do.Provide(injector, NewPublisher)
	do.Provide(injector, NewOrderSuccessConsumer)

	// Base infrastructure
	do.Provide(injector, NewBaseEventListener)

	return injector
}

func NewCampaignRepository(injector do.Injector) (secondary.CampaignRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](injector)
	return campaignPostgres.NewCampaignPersistenceRepository(getDbFunc), nil
}

func NewCampaignService(injector do.Injector) (primary.CampaignService, error) {
	campaignRepository := do.MustInvoke[secondary.CampaignRepository](injector)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](injector)
	return campaignService.NewCampaignService(
		campaignRepository,
		atomicExecutor,
	), nil
}

func NewOrderRepository(injector do.Injector) (secondary.OrderRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](injector)
	return orderPostgres.NewOrderPersistenceRepository(
		getDbFunc,
	), nil
}

func NewOrderService(injector do.Injector) (primary.OrderService, error) {
	orderRepository := do.MustInvoke[secondary.OrderRepository](injector)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](injector)
	return orderService.NewOrderService(
		orderRepository,
		atomicExecutor,
	), nil
}

func NewCampaignHandler(injector do.Injector) (campaignHandler.CampaignHandler, error) {
	service := do.MustInvoke[primary.CampaignService](injector)
	return campaignHandler.NewCampaignHandler(service), nil
}

func NewPublisher(injector do.Injector) (messagequeue.Publisher, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	return messagequeue.NewPublisher(cfg, tasks), nil
}

func NewBaseEventListener(injector do.Injector) (*messagequeue.BaseEventListener, error) {
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	logger := do.MustInvoke[*slog.Logger](injector)
	return messagequeue.NewBaseEventListener(tasks, logger), nil
}

func NewOrderSuccessConsumer(injector do.Injector) (*orderConsumer.OrderSuccessConsumer, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	orderService := do.MustInvoke[primary.OrderService](injector)
	baseEventListener := do.MustInvoke[*messagequeue.BaseEventListener](injector)
	return orderConsumer.NewOrderSuccessConsumer(baseEventListener, cfg.OrderSuccess, orderService), nil
}
