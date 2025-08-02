package fi_frontend

import (
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"specommerce/paymentservice/assets"
	"specommerce/paymentservice/config"
	"specommerce/paymentservice/di"
	paymentConsumer "specommerce/paymentservice/internal/adapters/primary/payment/event/kafka"
	"specommerce/paymentservice/pkg/atomicity"
	"specommerce/paymentservice/pkg/database"
	"specommerce/paymentservice/pkg/environment"
	"specommerce/paymentservice/pkg/service_config"
	"specommerce/paymentservice/pkg/shutdown"
	"specommerce/paymentservice/server"
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

	processPaymentRequestConsumer := do.MustInvoke[*paymentConsumer.ProcessPaymentRequestConsumer](injector)

	eg.Go(func() error {
		return processPaymentRequestConsumer.Start()
	})

	return eg.Wait()
}
