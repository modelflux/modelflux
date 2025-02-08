package workflow

import (
	"fmt"
	"regexp"

	generate "github.com/modelflux/cli/pkg/ai"
	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/tool"
)

// A WorkflowNode represents a step in the workflow.
type WorkflowNode struct {
	StepName  string
	ID        string
	Params    map[string]interface{}
	Operation string
	Tool      tool.Tool
	Model     model.Model
	Next      string // The ID of the next node to run
	Output    string
	Log       bool
}

func (n *WorkflowNode) replacePlaceholders(outputs map[string]string) {
	// Replace placeholders in the parameters
	for k, v := range n.Params {
		if s, ok := v.(string); ok {
			// Identify a placeholder in the string. ${{stepid.output}}
			regexp := regexp.MustCompile(`\${{([a-zA-Z0-9-]+)\.output}}`)
			matches := regexp.FindAllStringSubmatch(s, -1)
			for _, match := range matches {
				// Replace the placeholder with the output value.
				if len(match) == 2 {
					if output, ok := outputs[match[1]]; ok {
						s = regexp.ReplaceAllString(s, output)
						n.Params[k] = s
					}
				}
			}
		}
	}
}

func (n *WorkflowNode) Run(outputs map[string]string) (string, error) {
	fmt.Println("Running step:", n.StepName)

	// Replace placeholders in the parameters
	n.replacePlaceholders(outputs)

	var output string
	var err error
	if n.Tool != nil {
		output, err = n.Tool.Run(n.Params)
	} else if n.Model != nil {
		if n.Operation == "generate" {
			output, err = generate.Run(n.Params, n.Model)
		}
	}

	if err != nil {
		return "", err
	}

	n.Output = output
	fmt.Println("Output:", n.Output)

	return n.Next, nil
}
