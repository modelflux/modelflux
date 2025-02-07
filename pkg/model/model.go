package model

import (
	"fmt"

	"github.com/modelflux/cli/pkg/tool"
)

// Model is the interface all models must implement.
//
// Models are considered a special type of tool that can have additional standard methods.
// Methods:
// - New: Initialize the model. This is for models that require initialization. Other it will just return nil.

type Model interface {
	tool.Tool
	Generate(input string) (string, error)
}

func GetModel(name string) (Model, error) {
	switch name {
	case "ollama":
		return &OllamaModel{}, nil
	case "azure-openai":
		return &AzureOpenAIModel{}, nil
	case "openai":
		return &OpenAIModel{}, nil
	default:
		return nil, fmt.Errorf("model %s not found", name)
	}
}
