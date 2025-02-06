package workflow

import (
	"fmt"

	"github.com/modelflux/cli/pkg/load"
	"github.com/spf13/viper"
)

// WorkflowBuilder provides a chainable API to build a Workflow.
type WorkflowBuilder struct {
	wf   *Workflow
	data *load.WorkflowSchema
	cfg  viper.Viper
	err  error
}

// BuildWorkflow starts the workflow building process.
func BuildWorkflow(workflowName string, cfg *viper.Viper) *WorkflowBuilder {
	return &WorkflowBuilder{wf: &Workflow{name: workflowName}}
}

// Parse parses the YAML into the workflow structure.
func (b *WorkflowBuilder) Load() *WorkflowBuilder {
	if b.err != nil {
		return b
	}
	b.data = load.Load(b.wf.name)
	return b
}

// Validate checks the workflow schema.
func (b *WorkflowBuilder) Validate() *WorkflowBuilder {
	for _, m := range b.data.Models {
		// Get the model identifier
		// Validate the model options
		fmt.Println("validated model:", m)

	}
	for _, tool := range b.data.Tools {
		// Get the tool identifier
		// Validate the tool options
		fmt.Println("validated tool:", tool)

		// Add the tool to the workflow
	}

	// Create a root node for the workflow

	for _, step := range b.data.Task.Steps {
		// Validate the step
		fmt.Println("validated step:", step)

		// Create a workflow node

		// Add the node to the workflow

		// Connect the node to its dependencies
	}

	return b
}

// Init initializes any tools needed for the workflow.
// Returns the Workflow instance.
func (b *WorkflowBuilder) Init() *WorkflowBuilder {
	if b.err != nil {
		return b
	}
	// TODO: Initialize your tools here.
	fmt.Println("Tools initialized")
	return b
}

func (b *WorkflowBuilder) Build() (*Workflow, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.wf, nil
}
