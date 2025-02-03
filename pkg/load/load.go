package load

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type WorkflowSchema struct {
	Name   string           `yaml:"name"`
	Models map[string]Model `yaml:"models"`
	Tools  map[string]Tool  `yaml:"tools"`
	Task   Task             `yaml:"task"`
}

type Model interface {
}

type OpenAIModel struct {
	Key      string `yaml:"key"`
	Endpoint string `yaml:"endpoint"`
}

type AzureOpenAIModel struct {
	Key        string `yaml:"key"`
	Endpoint   string `yaml:"endpoint"`
	Version    string `yaml:"version"`
	Deployment string `yaml:"deployment"`
}

type Tool struct {
	Key         string            `yaml:"key"`
	Source      string            `yaml:"source"`
	Type        string            `yaml:"type"`
	ToolOptions map[string]string `yaml:"toolOptions"`
}

type Task struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name   string `yaml:"name"`
	Model  string `yaml:"model,omitempty"`
	Tool   string `yaml:"tool,omitempty"`
	Input  string `yaml:"input,omitempty"`
	Output string `yaml:"output,omitempty"`
}

func loadYAML(filePath string) (*WorkflowSchema, error) {
	data, err := os.ReadFile(filePath)
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

// ResolvePlaceholders resolves placeholders in task parameters
func ResolvePlaceholders(workflow *WorkflowSchema) {
	for j, step := range workflow.Task.Steps {
		if step.Input != "" && strings.Contains(step.Input, "{{") && strings.Contains(step.Input, "}}") {
			// Example: "{{task.steps.step1.output}}"
			// Resolve the placeholder here
			// This is a simplified example, you would need to implement the actual logic
			workflow.Task.Steps[j].Input = resolvePlaceholder(workflow, step.Input)
		}
	}
}

// resolvePlaceholder is a helper function to resolve a placeholder
func resolvePlaceholder(workflow *WorkflowSchema, placeholder string) string {
	// This is a simplified example, you would need to implement the actual logic
	// to parse the placeholder and retrieve the correct output value
	// Example: "{{task.steps.step1.output}}"
	parts := strings.Split(placeholder, ".")
	if len(parts) == 4 && parts[0] == "{{task" && parts[3] == "output}}" {
		stepName := parts[2]
		for _, step := range workflow.Task.Steps {
			if step.Name == stepName {
				return step.Output
			}
		}
	}
	return ""
}

func Load(workflowName string) *WorkflowSchema {
	workflowPath := fmt.Sprintf("workflows/%s.yaml", workflowName)
	fmt.Println("Loading workflow:", workflowPath)
	workflow, err := loadYAML(workflowPath)
	if err != nil {
		fmt.Println("Error loading YAML:", err)
		return nil
	}
	ResolvePlaceholders(workflow)
	return workflow
}
