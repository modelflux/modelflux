package model

import (
	"fmt"

	"github.com/spf13/viper"
)

type ModelConfiguration struct {
	Identifier   string                 `yaml:"identifier"`
	ModelOptions map[string]interface{} `yaml:"options,omitempty"`
}

// Model represents an interface for a model with methods to generate output,
// initialize the model, and validate and set options and parameters.
//
// Generate takes an input string and returns a generated string and an error if any occurs.
//
// New initializes the model and returns an error if any occurs.
//
// ValidateAndSetOptions validates and sets the provided user options using the given viper configuration,
// and returns an error if any occurs.
//
// ValidateAndSetParameters validates and sets the provided user parameters,
// and returns an error if any occurs.
type Model interface {
	Generate(input string) (string, error)
	New() error
	ValidateAndSetOptions(uOptions map[string]interface{}, cfg *viper.Viper) error
	ValidateAndSetParameters(uParams map[string]interface{}) error
}

func ValidateAndBuild(mdata ModelConfiguration, cfg *viper.Viper) (Model, error) {
	var m Model
	switch mdata.Identifier {
	case "openai":
		m = &OpenAIModel{}
	case "azure":
		m = &AzureOpenAIModel{}
	case "ollama":
		m = &OllamaModel{}
	default:
		return nil, fmt.Errorf("model %s not found", mdata.Identifier)
	}

	if err := m.ValidateAndSetOptions(mdata.ModelOptions, cfg); err != nil {
		return nil, fmt.Errorf("error validating model options: %v", err)
	}

	return m, nil
}
