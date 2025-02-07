package workflow

import (
	"fmt"
	"os"

	"github.com/modelflux/cli/pkg/tool"
	"gopkg.in/yaml.v3"
)

type WorkflowSchema struct {
	Models map[string]tool.ToolConfiguration `yaml:"models"`
	Tools  map[string]tool.ToolConfiguration `yaml:"tools"`
	Task   Task                              `yaml:"task"`
}

type Task struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name       string                 `yaml:"name"`
	ID         string                 `yaml:"id,omitempty"`
	Model      string                 `yaml:"model,omitempty"`
	Tool       string                 `yaml:"tool,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	Output     string                 `yaml:"output,omitempty"`
}

func LoadSchema(workflowName string) (*WorkflowSchema, error) {
	workflowPath := fmt.Sprintf("workflows/%s.yaml", workflowName)
	fmt.Println("Loading workflow:", workflowName)
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
