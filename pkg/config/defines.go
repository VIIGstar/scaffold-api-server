package config

// LogConfig holds details necessary for logging.
type LogConfig struct {
	// Format specifies the output log format.
	// Accepted values are: json, logfmt
	Format string

	// Level is the minimum log level that should appear on the output.
	Level string

	// NoColor makes sure that no log output gets colorized.
	NoColor bool
}

// Configuration holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type configuration struct {
	// Log configuration
	Log LogConfig

	// Telemetry configuration
	Telemetry struct {
		// Telemetry HTTP server address
		Addr string
	}

	// App configuration
	App appConfig

	// Database connection information
	Database DBConfig
}
