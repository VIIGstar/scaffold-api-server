package database

import (
	"fmt"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"gorm.io/gorm"
	"logur.dev/logur"
	"scaffold-api-server/internal/services"
	mysql "scaffold-api-server/internal/services/database/mysql"
	health_check "scaffold-api-server/internal/services/health-check"
	"scaffold-api-server/pkg"
	"scaffold-api-server/pkg/config"
	"time"
)

type DB struct {
	gormDB *gorm.DB
}

func New(logger logur.LoggerFacade, config config.DBConfig) *DB {
	healthListener := health_check.NewHealthListener(logger, "mysql")
	healthChecker := health.New(health.WithHealthListeners(healthListener))
	// Connect to the database
	logger.Info("connecting to database")
	mysql.SetLogger(logger)

	db, err := mysql.NewConnector(config)
	if err != nil {
		panic(fmt.Sprintf("connect database failed, error: %v", err))
	}
	sqlDB, _ := db.DB()

	// Register database health check
	_ = healthChecker.RegisterCheck(
		checks.Must(checks.NewPingCheck("db.check", sqlDB)),
		health.ExecutionPeriod(time.Minute*2))

	services.RegisterApp(string(pkg.MysqlConnectorAppName), db)
	return &DB{
		gormDB: db,
	}
}

func TestDB(db *gorm.DB) *DB {
	return &DB{
		gormDB: db,
	}
}

func (db *DB) GormDB() *gorm.DB {
	return db.gormDB
}

func (db *DB) Close() {
	dbConnection, _ := db.gormDB.DB()
	if dbConnection != nil {
		dbConnection.Close()
	}
}
