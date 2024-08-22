package cli

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetAPIKey() string {
	return os.Getenv("API_KEY")
}

func GetAzureOpenAIEndpoint() string {
	return os.Getenv("AZURE_OPENAI_ENDPOINT")
}
