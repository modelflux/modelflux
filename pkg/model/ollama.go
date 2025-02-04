package model

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

type OllamaModel struct {
	Model string `yaml:"model"`
}

func (o *OllamaModel) New() error {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	req := &api.PullRequest{
		Model: o.Model,
	}
	progressFunc := func(resp api.ProgressResponse) error {
		fmt.Printf("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
		return nil
	}

	err = client.Pull(ctx, req, progressFunc)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// PartialGenerateRequest contains only the fields needed for the Generate request.
type PartialGenerateRequest struct {
	Prompt string `json:"prompt"`
}

func (o *OllamaModel) Generate(input string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	req := &api.GenerateRequest{
		Model:  o.Model,
		Prompt: input,
	}

	ctx := context.Background()

	// Save the response to a variable
	var response string
	respFunc := func(resp api.GenerateResponse) error {
		fmt.Println(resp.Response)
		response = resp.Response
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	return response, nil
}
