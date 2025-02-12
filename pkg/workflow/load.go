package workflow

import (
	"fmt"
	"os"
	"path"

	"github.com/modelflux/cli/pkg/model"
	"gopkg.in/yaml.v3"
)

type WorkflowSchema struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	ID    string                   `yaml:"id,omitempty"`
	Name  string                   `yaml:"name"`
	Uses  string                   `yaml:"uses,omitempty"`
	Model model.ModelConfiguration `yaml:"model,omitempty"`
	Run   string                   `yaml:"run,omitempty"`  // Run is an operation to be preformed by the model
	With  map[string]interface{}   `yaml:"with,omitempty"` // With is the parameters to be passed to the tool.
	Log   bool                     `yaml:"log,omitempty"`  // Log is a flag wether the output of the tool should be logged to the console.
}

func LoadSchema(workflowName string, local bool) (*WorkflowSchema, error) {
	fmt.Println("LOADING WORKFLOW:", workflowName)

	var workflowPath string
	if local {
		workflowPath = workflowName + ".yaml"
		workflowPath = path.Clean(workflowPath)
	} else {
		workflowsDir := path.Join(".modelflux", "workflows")
		workflowFile := workflowName + ".yaml"
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		workflowPath = path.Join(home, workflowsDir, workflowFile)
		workflowPath = path.Clean(workflowPath)
	}

	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return nil, err
	}

	var workflow WorkflowSchema
	err = yaml.Unmarshal(data, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}
