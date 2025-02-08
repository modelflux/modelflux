package workflow

import (
	"fmt"
)

// Workflow represents your parsed and built workflow.
type Workflow struct {
	name     string
	graph    map[string]*WorkflowNode
	outputs  map[string]string
	rootNode string
}

// Initializes the models used in the workflow.
func (wf *Workflow) Init() error {
	for _, step := range wf.graph {
		if step.Model != nil {
			err := step.Model.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (wf *Workflow) Run() error {
	fmt.Println()
	fmt.Println("Running workflow", wf.name)
	// Starting at the root node, run each node in the graph
	// until there are no more nodes to run.
	n := wf.rootNode
	var err error
	for n != "" {
		fmt.Println("--------------------")
		node := wf.graph[n]
		n, err = node.Run(wf.outputs)
		if err != nil {
			return err
		}
		wf.outputs[node.ID] = node.Output
		fmt.Println("--------------------")
	}
	fmt.Println(("Workflow complete"))
	return nil
}
