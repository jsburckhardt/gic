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

type Config struct {
	LLMInstructions string `mapstructure:"llm_instructions"`
	ShouldCommit    bool   `mapstructure:"should_commit"`
	PR              bool   `mapstructure:"pr"`
}

type ConnectionConfig struct {
	ServiceProvider           string
	OpenAIAPIKey              string
	OpenAIAPIBase             string
	AzureAuthenticationType   string
	AzureOpenAIAPIKey         string
	AzureOpenAIEndpoint       string
	AzureOpenAIDeploymentName string
	OllamaAPIKey              string
	OllamaAPIBase             string
}

func LoadConfig() (Config, ConnectionConfig, error) {
	l := logger.GetLogger()
	var cfg Config
	var connCfg ConnectionConfig

	l.Debug("Current working directory: " + os.Getenv("PWD"))
	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	l.Debug("reading config from: " + os.Getenv("PWD") + "/.gic.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return cfg, connCfg, err
	}
	l.Debug("config file read successfully")
	l.Debug("unmarshalling config")
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, connCfg, err
	}
	l.Debug("config unmarshalled successfully")
	l.Debug("loading connection config from environment variables")
	connCfg = loadConnectionConfigFromEnv()
	l.Debug("validating config")
	if err := validateConfig(cfg, connCfg); err != nil {
		return cfg, connCfg, err
	}
	l.Debug("config validated successfully")
	return cfg, connCfg, nil
}

func loadConnectionConfigFromEnv() ConnectionConfig {
	return ConnectionConfig{
		ServiceProvider:           os.Getenv("SERVICE_PROVIDER"),
		OpenAIAPIKey:              os.Getenv("OPENAI_API_KEY"),
		OpenAIAPIBase:             os.Getenv("OPENAI_API_BASE"),
		AzureAuthenticationType:   os.Getenv("AZURE_AUTHENTICATION_TYPE"),
		AzureOpenAIAPIKey:         os.Getenv("AZURE_OPENAI_API_KEY"),
		AzureOpenAIEndpoint:       os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIDeploymentName: os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME"),
		OllamaAPIKey:              os.Getenv("OLLAMA_API_KEY"),
		OllamaAPIBase:             os.Getenv("OLLAMA_API_BASE"),
	}
}

func validateConfig(cfg Config, connCfg ConnectionConfig) error {
	l := logger.GetLogger()
	if cfg.LLMInstructions == emptyString {
		l.Debug("LLMInstructions not set in config. Using default value." + defaultInstructions)
		cfg.LLMInstructions = defaultInstructions
	}

	if err := validateConnectionConfig(connCfg); err != nil {
		return err
	}

	return nil
}

func validateConnectionConfig(connCfg ConnectionConfig) error {
	switch connCfg.ServiceProvider {
	case "openai":
		if connCfg.OpenAIAPIKey == emptyString {
			return fmt.Errorf("OPENAI_API_KEY environment variable not set")
		}
		if connCfg.OpenAIAPIBase == emptyString {
			return fmt.Errorf("OPENAI_API_BASE environment variable not set")
		}
	case "azure":
		if connCfg.AzureAuthenticationType == emptyString {
			return fmt.Errorf("AZURE_AUTHENTICATION_TYPE environment variable not set")
		}
		if connCfg.AzureAuthenticationType != "api_key" && connCfg.AzureAuthenticationType != "azure_ad" {
			return fmt.Errorf("AZURE_AUTHENTICATION_TYPE must be either 'api_key' or 'azure_ad'")
		}
		if connCfg.AzureAuthenticationType == "api_key" && connCfg.AzureOpenAIAPIKey == emptyString {
			return fmt.Errorf("AZURE_OPENAI_API_KEY environment variable not set for api_key authentication type")
		}
		if connCfg.AzureOpenAIEndpoint == emptyString {
			return fmt.Errorf("AZURE_OPENAI_ENDPOINT environment variable not set")
		}
		if connCfg.AzureOpenAIDeploymentName == emptyString {
			return fmt.Errorf("AZURE_OPENAI_DEPLOYMENT_NAME environment variable not set")
		}
	case "ollama":
		if connCfg.OllamaAPIKey == emptyString {
			return fmt.Errorf("OLLAMA_API_KEY environment variable not set")
		}
		if connCfg.OllamaAPIBase == emptyString {
			return fmt.Errorf("OLLAMA_API_BASE environment variable not set")
		}
	default:
		return fmt.Errorf("unsupported service provider: %s", connCfg.ServiceProvider)
	}
	return nil
}

func CreateSampleConfig() error {
	l := logger.GetLogger()
	l.Debug("Creating sample configuration")
	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.Set("llm_instructions",
		"You are a helpful assistant, that helps generating commit messages based on git diffs.",
	)
	viper.Set("should_commit", false)
	if err := viper.WriteConfigAs(".gic.yaml"); err != nil {
		return err
	}
	l.Debug("Sample configuration created successfully")
	return nil
}
