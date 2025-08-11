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
	"specommerce/campaignservice/pkg/cache"
	"specommerce/campaignservice/pkg/database"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/shutdown"
)

func NewInjector() do.Injector {
	injector := do.New()

	do.Provide(injector, NewCampaignRepository)
	do.Provide(injector, NewCampaignService)
	do.Provide(injector, NewCampaignHandler)

	do.Provide(injector, NewOrderRepository)
	do.Provide(injector, NewOrderService)

	do.Provide(injector, NewPublisher)
	do.Provide(injector, NewOrderConsumer)

	do.Provide(injector, NewBaseEventListener)
	do.Provide(injector, NewRedisClient)

	return injector
}

func NewCampaignRepository(injector do.Injector) (secondary.CampaignRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](injector)
	return campaignPostgres.NewCampaignPersistenceRepository(getDbFunc), nil
}

func NewCampaignService(injector do.Injector) (primary.CampaignService, error) {
	campaignRepository := do.MustInvoke[secondary.CampaignRepository](injector)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	cacheClient := do.MustInvoke[cache.Cache](injector)
	return campaignService.NewCampaignService(
		campaignRepository,
		atomicExecutor,
		cfg,
		cacheClient,
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
	cacheClient := do.MustInvoke[cache.Cache](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	return orderService.NewOrderService(
		orderRepository,
		atomicExecutor,
		cacheClient,
		cfg,
	), nil
}

func NewCampaignHandler(injector do.Injector) (campaignHandler.CampaignHandler, error) {
	service := do.MustInvoke[primary.CampaignService](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	return campaignHandler.NewCampaignHandler(service, cfg), nil
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

func NewOrderConsumer(injector do.Injector) (*orderConsumer.OrderConsumer, error) {
	cfg := do.MustInvoke[config.AppConfig](injector)
	orderService := do.MustInvoke[primary.OrderService](injector)
	baseEventListener := do.MustInvoke[*messagequeue.BaseEventListener](injector)
	return orderConsumer.NewOrderConsumer(baseEventListener, cfg.OrderSuccess, orderService), nil
}

func NewRedisClient(injector do.Injector) (cache.Cache, error) {
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	return cache.NewRedisClient(cfg.Redis, tasks), nil
}
