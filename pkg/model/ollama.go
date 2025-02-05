package model

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/spf13/viper"
)

type OllamaModel struct {
	Model string
}

func (o *OllamaModel) New(cfg *viper.Viper) error {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return err
	}

	ctx := context.Background()
	req := &api.PullRequest{
		Model: o.Model,
	}

	fmt.Printf("Checking if model %s is downloaded: ", o.Model)

	isDownloaded := false

	resp, err := client.List(ctx)

	if err == nil {
		for _, model := range resp.Models {
			if model.Name == o.Model {
				isDownloaded = true
				break
			}
		}
	}

	if isDownloaded {
		fmt.Println("Model already downloaded")
		return nil
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

func (o *OllamaModel) Generate(input string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	req := &api.GenerateRequest{
		Model:  o.Model,
		Prompt: input,

		// Set stream to false to get a single response
		Stream: new(bool),
	}

	ctx := context.TODO()

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
