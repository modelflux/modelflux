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

	// Pretty print workflow
	fmt.Printf("%+v\n", workflow)

	// TODO: Validate the workflow schema

	for _, step := range workflow.Task.Steps {
		fmt.Println("Running step:", step.Name)

		var m model.Model

		if step.Model != "" {
			// Load the model if the model type is ollama
			modelType := workflow.Models[step.Model].Type

			if modelType == "ollama" {
				m = &model.OllamaModel{}
				m.New()
			}
		}

		prompt := step.Parameters["prompt"]
		// Execute the prompt
		if step.Parameters["prompt"] != "" {
			fmt.Println("Prompt:", step.Parameters["prompt"])
			model.Generate(models.PartialGenerateRequest{Prompt: prompt})
		}
	}
}
