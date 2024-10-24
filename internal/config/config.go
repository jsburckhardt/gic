// Package config provides functionality for managing configuration settings.
package config

import (
	"fmt"
	"gic/internal/logger"
	"os"

	"github.com/spf13/viper"
)

const emptyString = ""
const zeroTokens = 0
const defaultTokens = 4000
const defaultInstructions = "You are a helpful assistant, that helps generating commit messages based on git diffs."

// Config represents the configuration structure for the application.
type Config struct {
	ModelDeploymentName string `mapstructure:"model_deployment_name"`
	APIVersion          string `mapstructure:"api_version"`
	LLMInstructions     string `mapstructure:"llm_instructions"`
	ConnectionType      string `mapstructure:"connection_type"`
	AzureEndpoint       string `mapstructure:"azure_endpoint"`
	Tokens              int    `mapstructure:"tokens"`
	ShouldCommit        bool   `mapstructure:"should_commit"`
	PR                  bool   `mapstructure:"pr"`
}

// LoadConfig loads the configuration from a YAML file and
// returns a Config struct.
func LoadConfig() (Config, error) {
	l := logger.GetLogger()
	var cfg Config

	l.Debug("Current working directory: " + os.Getenv("PWD"))
	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// print the path in which it is running
	l.Debug("reading config from: " + os.Getenv("PWD") + "/.gic.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	l.Debug("config file read successfully")
	l.Debug("unmarshalling config")
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	l.Debug("config unmarshalled successfully")
	l.Debug("validating config")
	if err := validateConfig(cfg); err != nil {
		return cfg, err
	}
	l.Debug("config validated successfully")
	return cfg, nil
}

// ValidateConfig validates the configuration and returns an
// error if any validation fails.
func validateConfig(cfg Config) error {
	l := logger.GetLogger()
	if cfg.LLMInstructions == emptyString {
		l.Debug("LLMInstructions not set in config. Using default value." + defaultInstructions)
		cfg.LLMInstructions = defaultInstructions
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

	return validateAPIVersion(cfg)
}

func validateAPIKey(cfg Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "openai" {
		if os.Getenv("API_KEY") == emptyString {
			return fmt.Errorf("API_KEY environment variable not set")
		}
	}
	return nil
}

func validateAzureEndpoint(cfg Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "azure_ad" {
		if cfg.AzureEndpoint == emptyString {
			return fmt.Errorf("AzureEndpoint not set in config")
		}
	}
	return nil
}

func validateTokens(cfg Config) error {
	if cfg.Tokens == zeroTokens {
		_, _ = fmt.Println("Tokens not set in config. Using default value 4000.")
		cfg.Tokens = defaultTokens
	}
	return nil
}

func validateModelDeploymentName(cfg Config) error {
	if cfg.ModelDeploymentName == emptyString {
		_, _ = fmt.Println("ModelDeploymentName not set in config. Using default value gpt-4o.")
		cfg.ModelDeploymentName = "gpt-4o"
	}
	return nil
}

func validateAPIVersion(cfg Config) error {
	if cfg.APIVersion == emptyString {
		_, _ = fmt.Println("ApiVersion not set in config. Using default value 2024-02-15-preview.")
		cfg.APIVersion = "2024-02-15-preview"
	}
	return nil
}

// CreateSampleConfig creates a sample configuration file.
func CreateSampleConfig() error {
	l := logger.GetLogger()
	l.Debug("Creating sample configuration")
	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.Set("model_deployment_name", "gpt-4o")
	viper.Set("api_version", "2024-02-15-preview")
	viper.Set("llm_instructions",
		"You are a helpful assistant, that helps generating commit messages based on git diffs.",
	)
	viper.Set("connection_type", "azure_ad")
	viper.Set("azure_endpoint", "https://<my_resource>.openai.azure.com/")
	viper.Set("should_commit", false)
	viper.Set("tokens", defaultTokens)
	if err := viper.WriteConfigAs(".gic.yaml"); err != nil {
		return err
	}
	l.Debug("Sample configuration created successfully")
	return nil
}
