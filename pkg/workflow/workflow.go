package workflow

import (
	"fmt"

	"github.com/modelflux/cli/pkg/load"
	"github.com/modelflux/cli/pkg/model"
)

type Tool struct {
}

type WorkflowNode[T Tool, M model.Model] struct {
	name       string
	input      []string // Identifiers of dependent nodes
	parameters interface{}
	tool       T
	model      M
	output     string // The result of running this node
	next       string
}

func (n *WorkflowNode[T, M]) Run() error {
	// Do somethings
	n.output = "output" + n.name
	return nil
}

// Workflow represents your parsed and built workflow.
type Workflow struct {
	name     string
	graph    map[string]*WorkflowNode[Tool, model.Model]
	rootNode string
}

// WorkflowBuilder provides a chainable API to build a Workflow.
type WorkflowBuilder struct {
	wf   *Workflow
	data *load.WorkflowSchema
	err  error
}

// BuildWorkflow starts the workflow building process.
func BuildWorkflow(workflowName string) *WorkflowBuilder {
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
	for _, model := range b.data.Models {
		// Get the model identifier
		// Validate the model options
		println("validated model:", model)

		// Add the model to the workflow
	}
	for _, tool := range b.data.Tools {
		// Get the tool identifier
		// Validate the tool options
		println("validated tool:", tool)

		// Add the tool to the workflow
	}

	// Create a root node for the workflow

	for _, step := range b.data.Task.Steps {
		// Validate the step
		println("validated step:", step)

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

func (wf *Workflow) Run() error {
	// Starting at the root node, run each node in the graph
	// until there are no more nodes to run.
	for node := wf.graph[wf.rootNode]; node != nil; node = wf.graph[node.next] {
		// Run the node
		err := node.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
