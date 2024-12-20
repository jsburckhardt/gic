// Package llm provides the logic for generating commit messages based on git diffs.
package llm

import (
	"context"
	"fmt"
	"log"

	"gic/internal/config"
	"gic/internal/logger"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/ollama/ollama/api"
)

const emptyString = ""
const responseMessage = 0

// GenerateCommitMessage generates a commit message based on the provided configuration and diff.
func GenerateCommitMessage(cfg config.Config, diff string) (string, error) {
	l := logger.GetLogger()
	l.Info("Generating commit message")
	if diff == emptyString {
		l.Info("No files staged for commit")
		return "### NO STAGED CHAGES ###", nil
	}

	switch cfg.ConnectionConfig.ServiceProvider {
	case "azure":
		return GenerateCommitMessageAzure(cfg, diff)
	case "openai":
		return GenerateCommitMessageOpenAI(cfg, diff)
	case "ollama":
		return GenerateCommitMessageOllama(cfg, diff)
	default:
		return emptyString, fmt.Errorf("unsupported connection type: %s", cfg.ConnectionConfig.ServiceProvider)
	}
}

// GenerateCommitMessageOllama generates a commit message using the Ollama service.
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
		Model:    cfg.ConnectionConfig.OllamaDeploymentName,
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

// GenerateCommitMessageAzure generates a commit message using the Azure service.
func GenerateCommitMessageAzure(cfg config.Config, diff string) (string, error) {
	var client *azopenai.Client
	var err error

	if cfg.ConnectionConfig.AzureAuthenticationType == "api_key" {
		keyCredential := azcore.NewKeyCredential(cfg.ConnectionConfig.AzureOpenAIAPIKey)
		client, err = azopenai.NewClientWithKeyCredential(cfg.ConnectionConfig.AzureOpenAIEndpoint, keyCredential, nil)
	} else if cfg.ConnectionConfig.AzureAuthenticationType == "azure_ad" {
		tokenCredential, tokenErr := azidentity.NewDefaultAzureCredential(nil)
		if tokenErr != nil {
			return emptyString, tokenErr
		}
		client, err = azopenai.NewClient(cfg.ConnectionConfig.AzureOpenAIEndpoint, tokenCredential, nil)
	} else {
		return emptyString, fmt.Errorf(
			"unsupported azure authentication type: %s",
			cfg.ConnectionConfig.AzureAuthenticationType,
		)
	}

	if err != nil {
		return emptyString, err
	}

	return getChatCompletions(cfg, client, diff)
}

// GenerateCommitMessageOpenAI generates a commit message using the OpenAI service.
func GenerateCommitMessageOpenAI(cfg config.Config, diff string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(cfg.ConnectionConfig.OpenAIAPIKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
		}),
		Model: openai.F(cfg.ConnectionConfig.OpenAIDeploymentName),
	})
	if err != nil {
		return emptyString, err
	}
	return chatCompletion.Choices[responseMessage].Message.Content, nil
}

func getChatCompletions(cfg config.Config, client *azopenai.Client, diff string) (string, error) {
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{

			Content: azopenai.NewChatRequestSystemMessageContent(cfg.LLMInstructions),
		},
		&azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent("git commit diff: " + diff),
		},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &(cfg.ConnectionConfig.AzureOpenAIDeploymentName),
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
