package log

import (
	"log"

	"logur.dev/logur"
)

// SetStandardLogger sets the global logger's output to a custom logger instance.
func SetStandardLogger(logger logur.Logger) {
	log.SetOutput(logur.NewLevelWriter(logger, logur.Info))
}
