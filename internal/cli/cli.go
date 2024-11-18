package cli

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func GetAPIKey() string {
	return os.Getenv("API_KEY")
}

func GetAzureOpenAIEndpoint() string {
	return os.Getenv("AZURE_OPENAI_ENDPOINT")
}

func GetModelDeploymentName() string {
	return os.Getenv("MODEL_DEPLOYMENT_NAME")
}

func GetAPIVersion() string {
	return os.Getenv("API_VERSION")
}

func GetConnectionType() string {
	return os.Getenv("CONNECTION_TYPE")
}

func GetAzureEndpoint() string {
	return os.Getenv("AZURE_ENDPOINT")
}

func GetTokens() int {
	tokensStr := os.Getenv("TOKENS")
	tokens, err := strconv.Atoi(tokensStr)
	if err != nil {
		return 4000 // default value
	}
	return tokens
}
