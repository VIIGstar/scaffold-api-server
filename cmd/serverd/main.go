package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"scaffold-api-server/cmd/serverd/router"
	"scaffold-api-server/internal/services/cache"
	"scaffold-api-server/internal/services/database"
	"scaffold-api-server/internal/services/log"
	"scaffold-api-server/pkg"
	build_info "scaffold-api-server/pkg/build-info"
	"scaffold-api-server/pkg/config"
	"syscall"
	"time"

	pkgerrors "github.com/pkg/errors"
	_ "scaffold-api-server/docs"
)

// Provisioned by ldflags
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

// @host localhost:4000
func run() error {
	ctx := context.TODO()
	v, f := viper.New(), pflag.NewFlagSet(string(pkg.APIAppName), pflag.ExitOnError)
	cfg := config.New(v, f)
	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(cfg.Log)

	// Override the global standard library logger to make sure everything uses our logger
	log.SetStandardLogger(logger)

	buildInfo := build_info.New(version, commitHash, buildDate)
	logger.Info("starting application", buildInfo.Fields())

	logger.Info("initialize database", buildInfo.Fields())
	db := database.New(logger, cfg.Database)
	defer db.Close()

	logger.Info("Initializing redis...")
	c, err := cache.New(ctx,
		fmt.Sprintf("%v:%v",
			viper.GetString("redis.host"),
			viper.GetString("redis.port"),
		),
		logger)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%s", cfg.App.HttpAddr)
	appEnv := viper.GetString("APPENV")
	server := &http.Server{
		Addr: addr,
		Handler: router.New(
			buildInfo,
			logger,
			db,
			c,
			appEnv == "dev",
			addr,
		),
	}

	// Create a channel to receive server error
	serverErr := make(chan error, 1)

	// Start on go routine for server
	go func() {
		logger.Info("Running server at: ", map[string]interface{}{
			"address": addr,
		})
		serverErr <- server.ListenAndServe()
	}()

	// Create a channel receive SIGTERM signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking code, waits for any signal to be received
	select {
	case <-shutdown:
		logger.Info("Shutting down server...")

		// Create context with timer for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		// Ask server to shutdown gracefully else force shutdown
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("could not stop server gracefully")
			return server.Close()
		}

		logger.Info("Shutdown server done!")

	case err := <-serverErr:
		return pkgerrors.WithStack(err)
	}
	return nil
}
