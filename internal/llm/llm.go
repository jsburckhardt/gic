package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/jsburckhardt/gic/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func GenerateCommitMessage(cfg config.Config, diff string) (string, error) {
	err := validateConfig(cfg)
	if err != nil {
		return "", err
	}
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
	return "", fmt.Errorf("not implemented")
}

func GenerateCommitMessageAzureAD(cfg config.Config, diff string) (string, error) {

	// const azureOpenAIEndpoint = "https://generic-aoai-01.openai.azure.com/"
	// The latest API versions, including previews, can be found here:
	// https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning
	// const azureOpenAIAPIVersion = "2024-02-15-preview"
	// modelDeploymentID := "gpt-4o"
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

func validateConfig(cfg config.Config) error {
	if cfg.LLMInstructions == "" {
		fmt.Println("LLMInstructions not set in config. Using default instructions.")
		cfg.LLMInstructions = "You are a helpful assistant, that helps generating commit messages based on git diffs."
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

	if err := validateApiVersion(cfg); err != nil {
		return err
	}

	return nil
}

func validateAPIKey(cfg config.Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "openai" {
		if os.Getenv("API_KEY") == "" {
			return fmt.Errorf("API_KEY environment variable not set")
		}
	}
	return nil
}

func validateAzureEndpoint(cfg config.Config) error {
	if cfg.ConnectionType == "azure" || cfg.ConnectionType == "azure_ad" {
		if cfg.AzureEndpoint == "" {
			return fmt.Errorf("AzureEndpoint not set in config")
		}
	}
	return nil
}

func validateTokens(cfg config.Config) error {
	if cfg.Tokens == 0 {
		fmt.Println("Tokens not set in config. Using default value 4000.")
		cfg.Tokens = 4000
	}
	return nil
}

func validateModelDeploymentName(cfg config.Config) error {
	if cfg.ModelDeploymentName == "" {
		fmt.Println("ModelDeploymentName not set in config. Using default value gpt-4o.")
		cfg.ModelDeploymentName = "gpt-4o"
	}
	return nil
}

func validateApiVersion(cfg config.Config) error {
	if cfg.ApiVersion == "" {
		fmt.Println("ApiVersion not set in config. Using default value 2024-02-15-preview.")
		cfg.ApiVersion = "2024-02-15-preview"
	}
	return nil
}
