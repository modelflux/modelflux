package model

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
)

type OpenAIModel struct {
	APIKey string
}

func (m *OpenAIModel) New(cfg *viper.Viper) error {
	mcfg := cfg.GetStringMapString("model")
	m.APIKey = mcfg["key"]

	if m.APIKey == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

func (m *OpenAIModel) Generate(input string) (string, error) {
	client := openai.NewClient(option.WithAPIKey(m.APIKey))
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
