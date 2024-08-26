// Package config provides functionality for managing configuration settings.
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the configuration structure for the application.
type Config struct {
	ModelDeploymentName string `mapstructure:"model_deployment_name"`
	APIVersion          string `mapstructure:"api_version"`
	LLMInstructions     string `mapstructure:"llm_instructions"`
	ConnectionType      string `mapstructure:"connection_type"`
	AzureEndpoint       string `mapstructure:"azure_endpoint"`
	Commit              bool   `mapstructure:"commit"`
	Tokens              int    `mapstructure:"tokens"`
}

// LoadConfig loads the configuration from a YAML file and returns a Config struct.
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

// ValidateConfig validates the configuration and returns an error if any validation fails.
func ValidateConfig(cfg Config) error {
	if cfg.LLMInstructions == "" {
		fmt.Println("LLMInstructions not set in config. Using default instructions.")
		cfg.LLMInstructions = "You are a helpful assistant, that helps generating commit messages based on git diffs."
	}

	if err := validateAPIKey(cfg); err != nil {
		return err
	}

	if err := validateAzureEndpoint(cfg); err != nil {
		return err
	}

	if err := validateTokens(cfg); err != nil {
		return err
	}

	if err := validateModelDeploymentName(cfg); err != nil {
		return err
	}

	if err := validateAPIVersion(cfg); err != nil {
		return err
	}

	return nil
}

func validateAPIKey(cfg Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "openai" {
		if os.Getenv("API_KEY") == "" {
			return fmt.Errorf("API_KEY environment variable not set")
		}
	}
	return nil
}

func validateAzureEndpoint(cfg Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "azure_ad" {
		if cfg.AzureEndpoint == "" {
			return fmt.Errorf("AzureEndpoint not set in config")
		}
	}
	return nil
}

func validateTokens(cfg Config) error {
	if cfg.Tokens == 0 {
		fmt.Println("Tokens not set in config. Using default value 4000.")
		cfg.Tokens = 4000
	}
	return nil
}

func validateModelDeploymentName(cfg Config) error {
	if cfg.ModelDeploymentName == "" {
		fmt.Println("ModelDeploymentName not set in config. Using default value gpt-4o.")
		cfg.ModelDeploymentName = "gpt-4o"
	}
	return nil
}

func validateAPIVersion(cfg Config) error {
	if cfg.APIVersion == "" {
		fmt.Println("ApiVersion not set in config. Using default value 2024-02-15-preview.")
		cfg.APIVersion = "2024-02-15-preview"
	}
	return nil
}
