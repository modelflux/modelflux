package run

import (
	"fmt"
	"log"
	"strings"

	"github.com/orelbn/tbd/pkg/load"
	"github.com/orelbn/tbd/pkg/model"
	"github.com/orelbn/tbd/pkg/tool"
)

func Run(workflowName string) {

	// The following code is for demonstration purposes only and will be refactored frequently.

	fmt.Println("run called")
	workflow := load.Load(workflowName)
	if workflow == nil {
		log.Fatal("Workflow loading failed.")
	}

	for _, step := range workflow.Task.Steps {
		// Print a new line between steps. This is done to make the output more readable.
		fmt.Println(" ")
		fmt.Println("Running step:", step.Name)

		toolCfg := workflow.Tools[step.Tool]
		modelCfg := workflow.Models[step.Model]

		if toolCfg.Identifier == "" && modelCfg.Identifier == "" {
			log.Fatal("Step must have either a tool or a model.")
		} else if toolCfg.Identifier != "" && modelCfg.Identifier != "" {
			log.Fatal("Step cannot have both a tool and a model.")
		}

		if toolCfg.Identifier != "" {
			switch toolCfg.Identifier {
			case "text-file-reader":
				fmt.Println("Running text-file-reader")
				tool := &tool.TextFileReaderTool{}
				params, err := tool.ValidateParams(step.Parameters)
				if err != nil {
					log.Fatal(err)
				}
				content, err := tool.Run(params)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Content of the file:")
				fmt.Println("---")
				// Print first 5 lines of the file.
				contentLines := strings.Split(content, "\n")
				for i := 0; i < 5 && i < len(contentLines); i++ {
					fmt.Println(contentLines[i])
				}
				fmt.Println("---")

			case "text-file-writer":
				fmt.Println("Running text-file-writer")
				tool := &tool.TextFileWriterTool{}
				params, err := tool.ValidateParams(step.Parameters)
				if err != nil {
					log.Fatal(err)
				}
				writtenPath, err := tool.Run(params)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("File written to:", writtenPath)
				println("")
			default:
				log.Fatalf("Unknown tool: %s", toolCfg.Identifier)
			}
		} else {
			switch modelCfg.Identifier {
			case "ollama":
				{
					fmt.Println("Running ollama")

					modelName, ok := modelCfg.ModelOptions["model"].(string)
					if !ok {
						log.Fatal("Model name not provided")
					}

					m := &model.OllamaModel{Model: modelName}
					err := m.New(nil)
					if err != nil {
						log.Fatal(err)
					}

					prompt := step.Parameters["prompt"].(string)
					output, err := m.Generate(prompt)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Generated text:")
					fmt.Println("---")
					fmt.Println(output)
					fmt.Println("---")
				}
			}
		}
	}
}
