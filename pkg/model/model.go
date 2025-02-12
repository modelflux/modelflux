package model

import (
	"fmt"

	"github.com/spf13/viper"
)

type ModelConfiguration struct {
	Provider string                 `yaml:"provider"`
	Options  map[string]interface{} `yaml:"options,omitempty"`
}

// Model is an interface that specifies the behavior required for any model implementation.
// It provides a standardized way to validate and apply options and parameters, initialize resources,
// and generate outputs based on the provided input.
//
// The interface includes the following methods:
//   - Init() error:
//     Initializes any required internal state or resources.
//   - Generate(input string) (string, error):
//     Processes the provided input and produces a generated output.
type Model interface {
	ValidateAndSetOptions(uOptions map[string]interface{}, cfg *viper.Viper) error
	Init() error
	Generate(input string) (string, error)
}

// TODO: refactor this as these are technically the model providers
func GetModel(name string) (Model, error) {
	switch name {
	case "ollama":
		return new(OllamaModel), nil
	case "azure-openai":
		return new(AzureOpenAIModel), nil
	case "openai":
		return new(OpenAIModel), nil
	default:
		return nil, fmt.Errorf("model %s not found", name)
	}
}
