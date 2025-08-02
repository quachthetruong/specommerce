package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/do/v2"
	"log"
	"log/slog"
	"net/http"
	"os"
	"specommerce/paymentservice/config"
	"specommerce/paymentservice/pkg/shutdown"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 1 * time.Second
)

func ServeHTTP(injector do.Injector) error {
	cfg := do.MustInvoke[config.AppConfig](injector)
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	logger := do.MustInvoke[*slog.Logger](injector)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes(injector),
		ErrorLog:     log.New(os.Stderr, "", 0),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	tasks.AddShutdownTask(
		func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, defaultShutdownPeriod)
			defer cancel()
			return srv.Shutdown(ctx)
		},
	)

	logger.Info(fmt.Sprintf("starting server on %s", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Info(fmt.Sprintf("stopped server on %s", srv.Addr))

	return nil
}
