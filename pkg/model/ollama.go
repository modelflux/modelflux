package model

import (
	"context"
	"fmt"
	"log"

	"github.com/modelflux/cli/pkg/util"
	"github.com/ollama/ollama/api"
	"github.com/spf13/viper"
)

type ollamaModelOptions struct {
	Model string `yaml:"model"`
}

type ollamaModelParameters struct {
	Input string `yaml:"input"`
}

type OllamaModel struct {
	options    ollamaModelOptions
	parameters ollamaModelParameters
}

func (o *OllamaModel) ValidateAndSetOptions(uOptions map[string]interface{}, cfg *viper.Viper) error {
	// Create a struct from the map using the util package.
	options, err := util.CreateStruct[ollamaModelOptions](uOptions)

	if err != nil {
		return err
	}

	if options.Model == "" {
		return fmt.Errorf("missing required option model for ollama model")
	}

	o.options = options

	return nil
}

func (o *OllamaModel) ValidateParameters(uParams map[string]interface{}) error {
	err := util.ValidateStructFields[ollamaModelParameters](uParams)
	return err
}

func (o *OllamaModel) SetParameters(params map[string]interface{}) error {
	p, err := util.CreateStruct[ollamaModelParameters](params)

	if err != nil {
		return err
	}

	// Additional checks.
	if p.Input == "" {
		return fmt.Errorf("missing required parameter: input in ollama model parameters")
	}

	o.parameters = p

	return nil
}

func (o *OllamaModel) New() error {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return err
	}
	ctx := context.Background()

	fmt.Printf("Checking if model %s is downloaded: ", o.options.Model)

	isDownloaded := false

	resp, err := client.List(ctx)

	if err == nil {
		for _, model := range resp.Models {
			if model.Name == o.options.Model {
				isDownloaded = true
				break
			}
		}
	}

	if isDownloaded {
		fmt.Println("Model already downloaded")
		return nil
	}

	req := &api.PullRequest{
		Model: o.options.Model,
	}

	fmt.Println("Model not downloaded, downloading now")
	progress := 0
	progressFunc := func(resp api.ProgressResponse) error {
		status := resp.Status
		// Handle total being 0
		if resp.Total == 0 {
			fmt.Printf("Status: %s\n", status)
		} else {
			val := int(resp.Completed * 100 / resp.Total)
			if val != progress && val%10 == 0 {
				fmt.Printf("Status: %s, Progress: %d%%\n", status, val)
				progress = val
			}
		}
		return nil
	}

	err = client.Pull(ctx, req, progressFunc)
	if err != nil {
		return err
	}
	return nil
}

func (o *OllamaModel) Run() (string, error) {
	return o.Generate(o.parameters.Input)
}

func (o *OllamaModel) Generate(input string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	req := &api.GenerateRequest{
		Model:  o.options.Model,
		Prompt: input,

		// Set stream to false to get a single response
		Stream: new(bool),
	}

	ctx := context.Background()

	response := ""
	respFunc := func(resp api.GenerateResponse) error {
		response = resp.Response
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	return response, nil
}
