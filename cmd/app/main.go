package main

import (
	"database/sql"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"logur.dev/logur"
	"os/signal"
	database "scaffold-api-server/cmd/services/database/mysql"
	health_check "scaffold-api-server/cmd/services/health-check"
	"scaffold-api-server/cmd/services/log"
	"syscall"
	"time"

	"scaffold-api-server/pkg/build-info"

	"fmt"
	"os"
)

// Provisioned by ldflags
var (
	version    string
	commitHash string
	buildDate  string
)

var logger logur.LoggerFacade

const (
	// appCode is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the application itself, the actual instance needs to be identified via environment
	// and other details.
	appCode = "scaffolding-api"

	// friendlyAppName is the visible name of the application.
	friendlyAppName = "Simple scaffolding api"
)

func main() {
	config := initConfig()
	// Create logger (first thing after configuration loading)
	logger = log.NewLogger(config.Log)

	// Override the global standard library logger to make sure everything uses our logger
	log.SetStandardLogger(logger)

	buildInfo := build_info.New(version, commitHash, buildDate)

	logger.Info("starting application", buildInfo.Fields())
	// Start database
	db := initDatabaseConnector(config)
	defer db.Close()
	// Start API server
	initAppServer(config)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL)
	<-ch
	shutdown()
}

func initConfig() configuration {
	v, f := viper.New(), pflag.NewFlagSet(friendlyAppName, pflag.ExitOnError)
	configure(v, f)

	f.String("config", "", "Configuration file")
	_ = f.Parse(os.Args[1:])

	if v, _ := f.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", friendlyAppName, version, commitHash, buildDate)

		os.Exit(0)
	}

	if c, _ := f.GetString("config"); c != "" {
		v.SetConfigFile(c)
	}

	err := v.ReadInConfig()
	if err != nil {
		panic("failed to read configuration")
	}
	notFoundErr, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if configFileNotFound {
		panic(notFoundErr.Error())
	}
	var config configuration
	err = v.Unmarshal(&config)
	return config
}

func initDatabaseConnector(config configuration) *sql.DB {
	healthListener := health_check.NewHealthListener(logger, "mysql")
	healthChecker := health.New(health.WithHealthListeners(healthListener))
	// Connect to the database
	logger.Info("connecting to database")
	dbConnector, err := database.NewConnector(config.Database)
	if err != nil {
		panic(fmt.Sprintf("connect database failed, error: %v", err))
	}

	database.SetLogger(logger)

	db := sql.OpenDB(dbConnector)

	// Register database health check
	_ = healthChecker.RegisterCheck(
		checks.Must(checks.NewPingCheck("db.check", db)),
		health.ExecutionPeriod(time.Second*10))
	return db
}

func shutdown() {
	logger.Info("shutshow server")
}
