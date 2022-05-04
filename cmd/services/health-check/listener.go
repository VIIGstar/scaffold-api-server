package health_check

import (
	health "github.com/AppsFlyer/go-sundheit"
	"logur.dev/logur"
)

type healthListener struct {
	name   string
	logger logur.Logger
}

func NewHealthListener(logger logur.Logger, serviceName string) health.HealthListener {
	return healthListener{
		name:   serviceName,
		logger: logger,
	}
}

func (c healthListener) OnResultsUpdated(results map[string]health.Result) {
	for _, result := range results {
		if result.Error != nil {
			c.logger.Info("health check failed", map[string]interface{}{"service": c.name, "error": result.Error.Error()})
			return
		}

		c.logger.Info("health check completed", map[string]interface{}{"service": c.name})
	}
}
