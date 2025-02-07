package model

import (
	"context"
	"fmt"

	"github.com/modelflux/cli/pkg/util"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/spf13/viper"
)

type azureOpenAIModelOptions struct {
	APIKey     string `yaml:"api_key"`
	Endpoint   string `yaml:"endpoint"`
	Deployment string `yaml:"deployment"`
	Version    string `yaml:"version"`
}

type azureOpenAIModelParameters struct {
	Input string `yaml:"input"`
}

type AzureOpenAIModel struct {
	options    azureOpenAIModelOptions
	parameters azureOpenAIModelParameters
}

func (m *AzureOpenAIModel) ValidateAndSetOptions(uOptions map[string]interface{}, cfg *viper.Viper) error {
	// Create a struct from the map using the util package.
	options, err := util.CreateStruct[azureOpenAIModelOptions](uOptions)

	if err != nil {
		return err
	}

	if options.APIKey == "" || options.Endpoint == "" || options.Deployment == "" || options.Version == "" {
		mcfg := cfg.GetStringMapString("model")
		m.options.APIKey = mcfg["key"]
		m.options.Endpoint = mcfg["endpoint"]
		m.options.Deployment = mcfg["deployment"]
		m.options.Version = mcfg["version"]
	} else {
		m.options = options
	}

	if m.options.APIKey == "" || m.options.Endpoint == "" || m.options.Deployment == "" || m.options.Version == "" {
		return fmt.Errorf("missing required api_key, endpoint, deployment, or version for azure model")
	}

	return nil
}

func (m *AzureOpenAIModel) ValidateParameters(uParams map[string]interface{}) error {
	err := util.ValidateStructFields[azureOpenAIModelParameters](uParams)
	return err
}

func (m *AzureOpenAIModel) SetParameters(params map[string]interface{}) error {
	p, err := util.CreateStruct[azureOpenAIModelParameters](params)

	if err != nil {
		return err
	}

	// Additional checks.
	if p.Input == "" {
		return fmt.Errorf("missing required parameter: input in azure-openai model parameters")
	}

	m.parameters = p

	return nil
}

func (m *AzureOpenAIModel) New() error {
	return nil
}

func (m *AzureOpenAIModel) Run() (string, error) {
	return m.Generate(m.parameters.Input)
}

func (m *AzureOpenAIModel) Generate(input string) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(m.options.Endpoint, m.options.Version),
		azure.WithAPIKey(m.options.APIKey),
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
