package run

import (
	"fmt"
	"log"

	"github.com/orelbn/tbd/pkg/load"
	"github.com/orelbn/tbd/pkg/model"
)

func Run(workflowName string) {
	fmt.Println("run called")
	workflow := load.Load(workflowName)
	if workflow == nil {
		log.Fatal("Workflow loading failed.")
	}

	// TODO: Validate the workflow schema

	for _, step := range workflow.Task.Steps {
		fmt.Println("Running step:", step.Name)

		var m model.Model

		if step.Model != "" {
			// Load the model if the model type is ollama
			modelType := workflow.Models[step.Model].Type

			if modelType == "ollama" {
				m = &model.OllamaModel{Model: workflow.Models[step.Model].Model}
				err := m.New(nil)
				if err != nil {
					fmt.Printf("failed to load model: %s", err)
					return
				}
				// Add space to output to make it easier to read
				fmt.Println(" ")
			}
		}

		prompt := step.Parameters["prompt"]
		// Execute the prompt
		if step.Parameters["prompt"] != "" {
			fmt.Println("Prompt:", step.Parameters["prompt"])
			response, err := m.Generate(prompt)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Response:", response)
			fmt.Println(" ")
		}
	}
}
