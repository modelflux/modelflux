package models

import (
	"context"

	"github.com/ollama/ollama/api"
)

type OllamaModel struct {
	Model string `yaml:"model"`
}

func (o *OllamaModel) Load() error {

	// Install ollama if it is not installed
	_, err := api.ClientFromEnvironment()
	if err != nil {
		// Ask the user if they want to install ollama
		// If yes, install ollama
		// If no, return an error

	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	err = client.Pull(context.Background(), &api.PullRequest{Model: o.Model}, nil)
	if err != nil {
		return err
	}

	client.Generate(context.Background(), &api.GenerateRequest{Model: o.Model}, nil)

	return nil
}
