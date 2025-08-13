package fi_frontend

import (
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"specommerce/campaignservice/assets"
	"specommerce/campaignservice/config"
	"specommerce/campaignservice/di"
	orderConsumer "specommerce/campaignservice/internal/adapters/primary/order/event/kafka"
	"specommerce/campaignservice/pkg/atomicity"
	"specommerce/campaignservice/pkg/database"
	"specommerce/campaignservice/pkg/environment"
	"specommerce/campaignservice/pkg/messagequeue"
	"specommerce/campaignservice/pkg/service_config"
	"specommerce/campaignservice/pkg/shutdown"
	"specommerce/campaignservice/server"
)

func Run(logger *slog.Logger, tasks *shutdown.Tasks) error {
	cfg, err := service_config.InitConfig[config.AppConfig](assets.EmbeddedFiles)
	if err != nil {
		return err
	}

	env := environment.Development
	if cfg.Env == string(environment.Production) {
		env = environment.Production
	}

	getDbFunc, atomicExecutor, err := database.New(cfg.Database, tasks, assets.EmbeddedFiles)
	if err != nil {
		return err
	}
	injector := di.NewInjector()
	do.ProvideValue(injector, logger)
	do.ProvideValue(injector, getDbFunc)
	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, env)
	do.ProvideValue[atomicity.AtomicExecutor](injector, atomicExecutor)
	do.ProvideValue(injector, tasks)

	var eg errgroup.Group
	eg.Go(
		func() error {
			return server.ServeHTTP(injector)
		})

	orderListener := do.MustInvoke[*orderConsumer.OrderConsumer](injector)
	successOrderListener := do.MustInvoke[*orderConsumer.SuccessOrderConsumer](injector)

	listeners := []messagequeue.EventListener{orderListener, successOrderListener}
	for _, l := range listeners {
		eg.Go(func() error {
			return l.Start()
		})
	}

	return eg.Wait()
}
