package config

import (
	"errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

// Process post-processes configuration after loading it.
func (configuration) Process() error {
	return nil
}

// Validate validates the configuration.
func (c configuration) Validate() error {
	if c.Telemetry.Addr == "" {
		return errors.New("telemetry http server address is required")
	}

	if err := c.App.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

// appConfig represents the application related configuration.
type appConfig struct {
	// HTTP server address
	// nolint: golint, stylecheck
	HttpAddr string

	// GRPC server address
	GrpcAddr string

	// Storage is the storage backend of the application
	Storage string
}

// Validate validates the configuration.
func (c appConfig) Validate() error {
	if c.HttpAddr == "" {
		return errors.New("http app server address is required")
	}

	if c.GrpcAddr == "" {
		return errors.New("grpc app server address is required")
	}

	if c.Storage != "inmemory" && c.Storage != "database" {
		return errors.New("app storage must be inmemory or database")
	}

	return nil
}

// New configures some defaults in the Viper instance.
func New(v *viper.Viper, f *pflag.FlagSet) configuration {
	bindDefault(v, f)
	f.String("config", "", "Configuration file")
	_ = f.Parse(os.Args[1:])

	if v, _ := f.GetBool("version"); v {
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

func bindDefault(v *viper.Viper, f *pflag.FlagSet) {
	v.SetConfigName("conf") // name of config file (without extension)
	v.SetConfigType("toml") // REQUIRED if the config file does not have the extension in the name
	// Viper settings
	v.AddConfigPath(".")
	//v.AddConfigPath("$CONFIG_DIR/")

	// Environment variable settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	// Global configuration
	v.SetDefault("shutdownTimeout", 15*time.Second)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		v.SetDefault("no_color", true)
	}

	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")
	v.RegisterAlias("log.noColor", "no_color")

	// App configuration
	f.String("http-addr", "8080", "App HTTP server address")
	_ = v.BindPFlag("app.httpAddr", f.Lookup("http-addr"))
	v.SetDefault("app.httpAddr", "8080")

	// Database configuration
	_ = v.BindEnv("database.host")
	v.SetDefault("database.port", 3306)
	_ = v.BindEnv("database.user")
	_ = v.BindEnv("database.pass")
	_ = v.BindEnv("database.name")
	v.SetDefault("database.params", map[string]string{
		"collation": "utf8mb4_general_ci",
	})
}
