package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"logur.dev/logur"
	"scaffold-api-server/internal/entities"
	"scaffold-api-server/internal/services"
	database "scaffold-api-server/internal/services/database/mysql"
	"scaffold-api-server/internal/services/log"
	"scaffold-api-server/pkg/config"
)

type MigrateService struct {
	services.DefaultService
	logger logur.LoggerFacade
	gormDB *gorm.DB
}

func main() {
	migrateService := MigrateService{}
	migrateService.Init()

	tables := []interface{}{
		entities.Investor{},
	}

	err := migrateService.gormDB.AutoMigrate(tables...)
	if err != nil {
		migrateService.logger.Error(fmt.Sprintf("Seed failed, details: %v", err))
		return
	}

	migrateService.logger.Info("Seed completed")
}

func (s *MigrateService) Init() {
	s.DefaultService.Init()
	var (
		logCf = config.LogConfig{}
		dbCf  = config.DBConfig{}
	)
	cfBytes, _ := json.Marshal(viper.GetStringMap("log"))
	json.Unmarshal(cfBytes, &logCf)
	cfBytes, _ = json.Marshal(viper.GetStringMap("database"))
	json.Unmarshal(cfBytes, &dbCf)

	logger := log.NewLogger(logCf)

	// Override the global standard library logger to make sure everything uses our logger
	log.SetStandardLogger(logger)
	// Start database
	if dbCf.Params == nil {
		dbCf.Params = make(map[string]string)
	}
	gormDB, err := database.NewConnector(dbCf)
	if err != nil {
		panic(err)
	}

	s.gormDB = gormDB
	s.logger = logger
}
