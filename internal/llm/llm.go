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

func GenerateCommitMessage(cfg config.Config, connCfg config.ConnectionConfig, diff string) (string, error) {
	l := logger.GetLogger()
	l.Info("Generating commit message")
	if diff == emptyString {
		l.Info("No files staged for commit")
		return "### NO STAGED CHAGES ###", nil
	}

	apikey := os.Getenv("API_KEY")
	switch connCfg.ConnectionType {
	case "azure":
		return GenerateCommitMessageAzure(apikey, connCfg, diff)
	case "azure_ad":
		return GenerateCommitMessageAzureAD(connCfg, diff, l)
	case "openai":
		return GenerateCommitMessageOpenAI(apikey, connCfg, diff)
	case "ollama":
		return GenerateCommitMessageOllama(connCfg, diff)
	default:
		return emptyString, fmt.Errorf("unsupported connection type: %s", connCfg.ConnectionType)
	}
}

func GenerateCommitMessageOllama(connCfg config.ConnectionConfig, diff string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return emptyString, err
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: connCfg.LLMInstructions,
		},
		{
			Role:    "user",
			Content: "git commit diff: " + diff,
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    connCfg.ModelDeploymentName,
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

func GenerateCommitMessageAzure(apikey string, connCfg config.ConnectionConfig, diff string) (string, error) {
	keyCredential := azcore.NewKeyCredential(apikey)
	client, err := azopenai.NewClientWithKeyCredential(connCfg.AzureEndpoint, keyCredential, nil)

	if err != nil {
		return emptyString, err
	}

	return getChatCompletions(connCfg, client, diff)
}

func GenerateCommitMessageAzureAD(connCfg config.ConnectionConfig, diff string, l *logger.Logger) (string, error) {
	l.Debug("GenerateCommitMessageAzureAD")
	l.Debug("obtaining token credential")
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		l.Error("failed getting token", "error", err)
		return emptyString, err
	}
	client, err := azopenai.NewClient(connCfg.AzureEndpoint, tokenCredential, nil)

	if err != nil {
		return emptyString, err
	}

	return getChatCompletions(connCfg, client, diff)
}

func GenerateCommitMessageOpenAI(apiKey string, connCfg config.ConnectionConfig, diff string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
		}),
		Model: openai.F(connCfg.ModelDeploymentName),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[responseMessage].Message.Content, nil
}

func getChatCompletions(connCfg config.ConnectionConfig, client *azopenai.Client, diff string) (string, error) {
	maxTokens := int32(connCfg.Tokens)

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{
			Content: to.Ptr(connCfg.LLMInstructions),
		},
		&azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent("git commit diff: " + diff),
		},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &(connCfg.ModelDeploymentName),
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
