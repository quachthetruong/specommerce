package fi_frontend

import (
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"specommerce/orderservice/assets"
	"specommerce/orderservice/config"
	"specommerce/orderservice/di"
	paymentConsumer "specommerce/orderservice/internal/adapters/primary/payment/event/kafka"
	"specommerce/orderservice/pkg/atomicity"
	"specommerce/orderservice/pkg/database"
	"specommerce/orderservice/pkg/environment"
	"specommerce/orderservice/pkg/service_config"
	"specommerce/orderservice/pkg/shutdown"
	"specommerce/orderservice/server"
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

	processPaymentResponseConsumer := do.MustInvoke[*paymentConsumer.ProcessPaymentResponseConsumer](injector)

	eg.Go(func() error {
		return processPaymentResponseConsumer.Start()
	})

	return eg.Wait()
}
