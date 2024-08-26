// Package cli provides a command-line interface for the application.
package cli

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
// It attempts to load the .env file located in the current working directory.
// If the file is not found or there is an error loading the file, it will cause a fatal error.
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetAPIKey returns the API key from the environment variables.
func GetAPIKey() string {
	return os.Getenv("API_KEY")
}

// GetAzureOpenAIEndpoint returns the Azure OpenAI endpoint from the environment variables.
func GetAzureOpenAIEndpoint() string {
	return os.Getenv("AZURE_OPENAI_ENDPOINT")
}
