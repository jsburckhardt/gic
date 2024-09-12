// Package llm provides functionality for interacting with the Language Learning Model.
package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"gic/internal/config"
	"gic/internal/logger"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/ollama/ollama/api"
)

const emptyString = ""
const responseMessage = 0

// GenerateCommitMessage generates a commit message based on the provided configuration and diff.
// It takes a config.Config object and a string representing the diff as input.
// The function returns a string containing the generated commit message and an error if any.
// The commit message is generated based on the connection type specified in the config.Config object.
// Supported connection types are "azure", "azure_ad", and "openai".
// If the connection type is not supported, the function returns an empty string and an error
// indicating the unsupported connection type.
func GenerateCommitMessage(cfg config.Config, diff string) (string, error) {
	l := logger.GetLogger()
	l.Info("Generating commit message")
	// if diff is empty finish
	if diff == emptyString {
		l.Debug("No files staged for commit")
		return "### NO STAGED CHAGES ###", nil
	}

	apikey := os.Getenv("API_KEY")
	switch cfg.ConnectionType {
	case "azure":
		return GenerateCommitMessageAzure(apikey, cfg, diff)
	case "azure_ad":
		return GenerateCommitMessageAzureAD(cfg, diff, l)
	case "openai":
		return GenerateCommitMessageOpenAI(apikey, cfg, diff)
	case "ollama":
		return GenerateCommitMessageOllama(cfg, diff)
	default:
		return emptyString, fmt.Errorf("unsupported connection type: %s", cfg.ConnectionType)
	}
}

// GenerateCommitMessageOllama generates a commit message using an LLM hosted in Ollama.
func GenerateCommitMessageOllama(cfg config.Config, diff string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return emptyString, err
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: cfg.LLMInstructions,
		},
		{
			Role:    "user",
			Content: "git commit diff: " + diff,
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    cfg.ModelDeploymentName,
		Messages: messages,
		Stream:   func(b bool) *bool { return &b }(false),
	}

	var commitMessage string
	respFunc := func(resp api.ChatResponse) error {
		commitMessage = resp.Message.Content
		return nil
	}
	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		return emptyString, err
	}
	return commitMessage, nil
}

// GenerateCommitMessageAzure generates a commit message using the Azure Language Learning Model.
// It takes an API key, a config.Config object, and a string representing the diff as input.
func GenerateCommitMessageAzure(apikey string, cfg config.Config, diff string) (string, error) {
	keyCredential := azcore.NewKeyCredential(apikey)
	client, err := azopenai.NewClientWithKeyCredential(cfg.AzureEndpoint, keyCredential, nil)

	if err != nil {
		return emptyString, err
	}

	return getChatCompletions(cfg, client, diff)
}

// GenerateCommitMessageAzureAD generates a commit message using
// the Azure Language Learning Model with Azure Active Directory
// authentication.
// It takes a config.Config object and a string representing
// the diff as input.
func GenerateCommitMessageAzureAD(cfg config.Config, diff string, l *logger.Logger) (string, error) {
	l.Debug("GenerateCommitMessageAzureAD")
	l.Debug("obtaining token credential")
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		l.Error("failed getting token", "error", err)
		return emptyString, err
	}
	client, err := azopenai.NewClient(cfg.AzureEndpoint, tokenCredential, nil)

	if err != nil {
		return emptyString, err
	}

	return getChatCompletions(cfg, client, diff)
}

// GenerateCommitMessageOpenAI generates a commit message using the OpenAI Language Learning Model.
// It takes an API key, a config.Config object, and a string representing the diff as input.
func GenerateCommitMessageOpenAI(apiKey string, cfg config.Config, diff string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
		}),
		Model: openai.F(cfg.ModelDeploymentName),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[responseMessage].Message.Content, nil
}

func getChatCompletions(cfg config.Config, client *azopenai.Client, diff string) (string, error) {
	maxTokens := int32(cfg.Tokens)

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{
			Content: to.Ptr(cfg.LLMInstructions),
		},
		&azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent("git commit diff: " + diff),
		},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &(cfg.ModelDeploymentName),
		MaxTokens:      &maxTokens,
	}, nil)

	if err != nil {
		log.Printf("ERROR: %s", err)
		return emptyString, err
	}

	var messageContent string
	for _, choice := range resp.Choices {
		if choice.ContentFilterResults != nil {
			if choice.ContentFilterResults.Error != nil {
				return emptyString, choice.ContentFilterResults.Error
			}
		}
		messageContent = *choice.Message.Content
	}

	return messageContent, nil
}
