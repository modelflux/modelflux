package workflow

import (
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
