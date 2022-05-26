package services

import (
	"github.com/spf13/viper"
)

var initiated = false

func init() {
	if !initiated {
		viper.SetConfigName("conf") // name of config file (without extension)
		viper.SetConfigType("toml") // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")    // optionally look for config in the working directory
		viper.ReadInConfig()
		initiated = true
	}
}

type DefaultService struct {
}

func (s *DefaultService) Init() {

}
