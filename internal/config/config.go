package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	LLMInstructions     string `mapstructure:"llm_instructions"`
	ConnectionType      string `mapstructure:"connection_type"`
	AzureEndpoint       string `mapstructure:"azure_endpoint"`
	Commit              bool   `mapstructure:"commit"`
	Tokens              int    `mapstructure:"tokens"`
	ModelDeploymentName string `mapstructure:"model_deployment_name"`
	ApiVersion          string `mapstructure:"api_version"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// print the path in which it is running

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
