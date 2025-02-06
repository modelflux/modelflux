package workflow

import "fmt"

// Initializes the models used in the workflow.
func (wf *Workflow) Init() {
	fmt.Println()
	// Initialize each model
	fmt.Println("Initializing models")
	for _, n := range wf.graph {
		if n.model != nil {
			(*n.model).New()
		}
	}
}

func (wf *Workflow) Run() error {
	fmt.Println()
	fmt.Println("Running workflow")

	fmt.Println("Task:", wf.task)
	// Starting at the root node, run each node in the graph
	// until there are no more nodes to run.
	n := wf.rootNode
	var err error
	for n != "" {
		node := wf.graph[n]
		n, err = node.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
