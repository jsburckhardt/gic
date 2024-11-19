package config

import (
	"fmt"
	"gic/internal/logger"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const emptyString = ""
const defaultInstructions = "You are a helpful assistant, that helps generating commit messages based on git diffs."
const defaultOpenAIDeploymentName = "gpt-4o-mini"
const defaultOllamaDeploymentName = "phi3"

type Config struct {
	LLMInstructions  string           `mapstructure:"llm_instructions"`
	ShouldCommit     bool             `mapstructure:"should_commit"`
	PR               bool             `mapstructure:"pr"`
	ConnectionConfig ConnectionConfig `mapstructure:"connection_config"`
}

type ConnectionConfig struct {
	ServiceProvider           string
	OpenAIAPIKey              string
	OpenAIAPIBase             string
	OpenAIDeploymentName      string
	AzureAuthenticationType   string
	AzureOpenAIAPIKey         string
	AzureOpenAIEndpoint       string
	AzureOpenAIDeploymentName string
	OllamaAPIKey              string
	OllamaAPIBase             string
	OllamaDeploymentName      string
}

func LoadConfig() (Config, error) {
	l := logger.GetLogger()
	var cfg Config

	l.Debug("Current working directory: " + os.Getenv("PWD"))
	viper.SetConfigName(".gic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

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
	l.Debug("loading connection config from environment variables")
	cfg.ConnectionConfig = loadConnectionConfigFromEnv()
	l.Debug("validating config")
	if err := validateConfig(cfg); err != nil {
		return cfg, err
	}
	l.Debug("config validated successfully")
	return cfg, nil
}

func loadConnectionConfigFromEnv() ConnectionConfig {
	l := logger.GetLogger()
	l.Debug("Loading connection config from .env file")
	if err := godotenv.Load(); err != nil {
		l.Warn("No .env file found or unable to load it. Using environment variables")
	}
	openAIDeploymentName := os.Getenv("OPENAI_DEPLOYMENT_NAME")
	if openAIDeploymentName == emptyString {
		openAIDeploymentName = defaultOpenAIDeploymentName
	}
	ollamaDeploymentName := os.Getenv("OLLAMA_DEPLOYMENT_NAME")
	if ollamaDeploymentName == emptyString {
		ollamaDeploymentName = defaultOllamaDeploymentName
	}
	return ConnectionConfig{
		ServiceProvider:           os.Getenv("SERVICE_PROVIDER"),
		OpenAIAPIKey:              os.Getenv("OPENAI_API_KEY"),
		OpenAIAPIBase:             os.Getenv("OPENAI_API_BASE"),
		OpenAIDeploymentName:      openAIDeploymentName,
		AzureAuthenticationType:   os.Getenv("AZURE_AUTHENTICATION_TYPE"),
		AzureOpenAIAPIKey:         os.Getenv("AZURE_OPENAI_API_KEY"),
		AzureOpenAIEndpoint:       os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIDeploymentName: os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME"),
		OllamaAPIKey:              os.Getenv("OLLAMA_API_KEY"),
		OllamaAPIBase:             os.Getenv("OLLAMA_API_BASE"),
		OllamaDeploymentName:      ollamaDeploymentName,
	}
}

func validateConfig(cfg Config) error {
	l := logger.GetLogger()
	l.Debug("Validating config")
	if cfg.LLMInstructions == emptyString {
		l.Debug("LLMInstructions not set in config. Using default value." + defaultInstructions)
		cfg.LLMInstructions = defaultInstructions
	}

	if err := validateConnectionConfig(cfg.ConnectionConfig); err != nil {
		return err
	}

	return nil
}

func validateConnectionConfig(connCfg ConnectionConfig) error {
	l := logger.GetLogger()
	l.Debug("Validating connection config from environment")
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
		if connCfg.OllamaDeploymentName == emptyString {
			return fmt.Errorf("OLLAMA_DEPLOYMENT_NAME environment variable not set")
		}
	default:
		return fmt.Errorf("unsupported service provider. options are openai, azure or ollama. got: %s. Options are openai, azure or ollama", connCfg.ServiceProvider)
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

func CreateSampleDotEnv() error {
	l := logger.GetLogger()
	l.Debug("Creating sample .env configuration")
	content := `SERVICE_PROVIDER=openai # openai, azure, ollama
OPENAI_API_KEY=your_openai_api_key # Required if SERVICE_PROVIDER=openai
OPENAI_API_BASE=https://api.openai.com/v1 # Required if SERVICE_PROVIDER=openai
OPENAI_DEPLOYMENT_NAME=gpt-4o-mini # Value for OpenAI deployment name defaults to gpt-4o-mini
AZURE_AUTHENTICATION_TYPE=api_key # api_key, azure_ad
AZURE_OPENAI_API_KEY=your_azure_openai_api_key # Required if SERVICE_PROVIDER=azure and AZURE_AUTHENTICATION_TYPE=api_key
AZURE_OPENAI_ENDPOINT=https://your-azure-endpoint # Required if SERVICE_PROVIDER=azure
AZURE_OPENAI_DEPLOYMENT_NAME=your-deployment-name # Required if SERVICE_PROVIDER=azure
OLLAMA_API_KEY=your_ollama_api_key # Required if SERVICE_PROVIDER=ollama
OLLAMA_API_BASE=https://api.ollama.com/v1 # Required if SERVICE_PROVIDER=ollama
OLLAMA_DEPLOYMENT_NAME=phi3 # Value for Ollama deployment name defaults to phi3`
	file, err := os.Create("sample.gic.env")
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	l.Debug(".env configuration created successfully")
	return nil
}
