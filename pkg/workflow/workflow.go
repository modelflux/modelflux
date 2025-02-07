package workflow

import (
	"fmt"

	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/tool"
)

// Workflow represents your parsed and built workflow.
type Workflow struct {
	Task     string
	Graph    map[string]*WorkflowNode
	Tools    map[string]*tool.Tool
	Models   map[string]*model.Model
	Outputs  map[string]string
	RootNode string
}

// Initializes the models used in the workflow.
func (wf *Workflow) Init() error {
	fmt.Println()
	// Initialize each model
	fmt.Println("Initializing models")
	for k, m := range wf.Models {
		err := (*m).New()
		if err != nil {
			return fmt.Errorf("error initializing model %s: %v", k, err)
		}
	}
	return nil
}

func (wf *Workflow) Run() error {
	fmt.Println()
	fmt.Println("Running workflow")

	fmt.Println("Task:", wf.Task)
	// Starting at the root node, run each node in the graph
	// until there are no more nodes to run.
	n := wf.RootNode
	var err error
	for n != "" {
		fmt.Println("--------------------")
		node := wf.Graph[n]
		n, err = node.Run(wf.Outputs)
		if err != nil {
			return err
		}
		wf.Outputs[node.ID] = node.Output
		fmt.Println("--------------------")
	}
	fmt.Println(("Workflow complete"))
	return nil
}
