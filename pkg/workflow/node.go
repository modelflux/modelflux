package workflow

import (
	"fmt"
	"regexp"

	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/tool"
)

type WorkflowNode struct {
	StepName   string
	ID         string
	Parameters map[string]interface{}
	Tool       *tool.Tool
	Model      *model.Model
	Output     string // The result of running this node
	Next       string
}

func (n *WorkflowNode) ReplacePlaceholders(outputs map[string]string) {
	// Replace placeholders in the parameters
	for k, v := range n.Parameters {
		if s, ok := v.(string); ok {
			// Identify a placeholder in the string. ${{stepid.output}}
			regexp := regexp.MustCompile(`\${{([a-zA-Z0-9-]+)\.output}}`)
			matches := regexp.FindAllStringSubmatch(s, -1)
			for _, match := range matches {
				// Replace the placeholder with the output value.
				if len(match) == 2 {
					if output, ok := outputs[match[1]]; ok {
						s = regexp.ReplaceAllString(s, output)
						n.Parameters[k] = s
					}
				}
			}
		}
	}
}

func (n *WorkflowNode) Run(outputs map[string]string) (string, error) {
	fmt.Println("Running step:", n.StepName)

	// Replace placeholders in the parameters
	n.ReplacePlaceholders(outputs)

	var output string
	var err error
	if n.Model != nil {
		(*n.Model).SetParameters(n.Parameters)
		output, err = (*n.Model).Run()
	} else if n.Tool != nil {
		(*n.Tool).SetParameters(n.Parameters)
		output, err = (*n.Tool).Run()
	}

	if err != nil {
		return "", err
	}

	n.Output = output
	fmt.Println("Output:", n.Output)

	return n.Next, nil
}
