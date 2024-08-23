package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/jsburckhardt/gic/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func GenerateCommitMessage(cfg config.Config, diff string) (string, error) {
	apikey := os.Getenv("API_KEY")
	switch cfg.ConnectionType {
	case "azure":
		return GenerateCommitMessageAzure(apikey, cfg, diff)
	case "azure_ad":
		return GenerateCommitMessageAzureAD(cfg, diff)
	case "openai":
		return GenerateCommitMessageOpenAI(apikey, cfg, diff)
	default:
		return "", fmt.Errorf("unsupported connection type: %s", cfg.ConnectionType)
	}
}

func GenerateCommitMessageAzure(apikey string, cfg config.Config, diff string) (string, error) {
	maxTokens := int32(cfg.Tokens)
	keyCredential := azcore.NewKeyCredential(apikey)

	client, err := azopenai.NewClientWithKeyCredential(cfg.AzureEndpoint, keyCredential, nil)

	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr(cfg.LLMInstructions)},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("git commit diff: " + diff)},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &(cfg.ModelDeploymentName),
		MaxTokens:      &maxTokens,
	}, nil)

	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}

	// var choice azopenai.ChatChoice
	var messageContent string
	for _, choice := range resp.Choices {
		if choice.ContentFilterResults != nil {

			if choice.ContentFilterResults.Error != nil {
				fmt.Fprintf(os.Stderr, "  Error:%v\n", choice.ContentFilterResults.Error)
			}

			// TODO: Include in logger
			// fmt.Fprintf(os.Stderr, "Content filter results\n")
			// fmt.Fprintf(os.Stderr, "  Hate: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Hate.Severity, *choice.ContentFilterResults.Hate.Filtered)
			// fmt.Fprintf(os.Stderr, "  SelfHarm: sev: %v, filtered: %v\n", *choice.ContentFilterResults.SelfHarm.Severity, *choice.ContentFilterResults.SelfHarm.Filtered)
			// fmt.Fprintf(os.Stderr, "  Sexual: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Sexual.Severity, *choice.ContentFilterResults.Sexual.Filtered)
			// fmt.Fprintf(os.Stderr, "  Violence: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Violence.Severity, *choice.ContentFilterResults.Violence.Filtered)
		}

		// TODO: Include in logger
		// if choice.Message != nil && choice.Message.Content != nil {
		// 	fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
		// }

		// TODO: Include in logger
		// if choice.FinishReason != nil {
		//// this choice's conversation is complete.
		// 	fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", *choice.Index, *choice.FinishReason)
		// }
		messageContent = *choice.Message.Content
	}

	return messageContent, nil
}

func GenerateCommitMessageAzureAD(cfg config.Config, diff string) (string, error) {
	maxTokens := int32(cfg.Tokens)
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("Failed to create the DefaultAzureCredential: %s", err)
		os.Exit(1)
	}

	client, err := azopenai.NewClient(cfg.AzureEndpoint, tokenCredential, nil)

	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}

	messages := []azopenai.ChatRequestMessageClassification{
		// You set the tone and rules of the conversation with a prompt as the system role.
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr(cfg.LLMInstructions)},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("git commit diff: " + diff)},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &(cfg.ModelDeploymentName),
		MaxTokens:      &maxTokens,
	}, nil)

	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}

	// var choice azopenai.ChatChoice
	var messageContent string
	for _, choice := range resp.Choices {
		if choice.ContentFilterResults != nil {

			if choice.ContentFilterResults.Error != nil {
				fmt.Fprintf(os.Stderr, "  Error:%v\n", choice.ContentFilterResults.Error)
			}

			// TODO: Include in logger
			// fmt.Fprintf(os.Stderr, "Content filter results\n")
			// fmt.Fprintf(os.Stderr, "  Hate: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Hate.Severity, *choice.ContentFilterResults.Hate.Filtered)
			// fmt.Fprintf(os.Stderr, "  SelfHarm: sev: %v, filtered: %v\n", *choice.ContentFilterResults.SelfHarm.Severity, *choice.ContentFilterResults.SelfHarm.Filtered)
			// fmt.Fprintf(os.Stderr, "  Sexual: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Sexual.Severity, *choice.ContentFilterResults.Sexual.Filtered)
			// fmt.Fprintf(os.Stderr, "  Violence: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Violence.Severity, *choice.ContentFilterResults.Violence.Filtered)
		}

		// TODO: Include in logger
		// if choice.Message != nil && choice.Message.Content != nil {
		// 	fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
		// }

		// TODO: Include in logger
		// if choice.FinishReason != nil {
		//// this choice's conversation is complete.
		// 	fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", *choice.Index, *choice.FinishReason)
		// }
		messageContent = *choice.Message.Content
	}

	return messageContent, nil
}

func GenerateCommitMessageOpenAI(apiKey string, cfg config.Config, diff string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
		}),
		Model: openai.F(openai.ChatModelGPT4),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
