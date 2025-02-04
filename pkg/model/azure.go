package model

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/spf13/viper"
)

type AzureOpenAIModel struct {
	APIKey     string
	Endpoint   string
	Deployment string
	Version    string
}

func (m *AzureOpenAIModel) New(cfg *viper.Viper) error {
	mcfg := cfg.GetStringMapString("model")
	m.APIKey = mcfg["key"]
	m.Endpoint = mcfg["endpoint"]
	m.Deployment = mcfg["deployment"]
	m.Version = mcfg["version"]

	if m.APIKey == "" || m.Endpoint == "" || m.Deployment == "" || m.Version == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

func (m *AzureOpenAIModel) Generate(input string) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(m.Endpoint, m.Version),
		azure.WithAPIKey(m.APIKey),
	)

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err.Error())
	}
	return resp.Choices[0].Message.Content, nil
}
