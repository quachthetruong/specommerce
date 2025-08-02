// @title Order Service API
// @version 1.0
// @description Order Service API for managing orders
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"github.com/shopspring/decimal"
	"log/slog"
	"os"
	"runtime/debug"
	app "specommerce/campaignservice"
	"specommerce/campaignservice/pkg/shutdown"
)

func main() {
	decimal.MarshalJSONWithoutQuotes = true
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tasks, _ := shutdown.NewShutdownTasks(logger)
	defer func() {
		tasks.Wait(recover())
	}()
	err := app.Run(logger, tasks)
	if err != nil {
		trace := debug.Stack()
		logger.Error("cannot start application", slog.String("error", err.Error()), slog.String("stack", string(trace)))
		os.Exit(1)
	}
}
