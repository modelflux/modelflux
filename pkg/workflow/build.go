package workflow

import (
	"fmt"

	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/tool"
	"github.com/modelflux/cli/pkg/util"
	"github.com/spf13/viper"
)

func ValidateAndBuildWorkflow(data *WorkflowSchema, cfg *viper.Viper) (*Workflow, error) {
	fmt.Println()

	// Initialize the workflow
	wf := &Workflow{}
	wf.Graph = make(map[string]*WorkflowNode)
	wf.Tools = make(map[string]*tool.Tool)
	wf.Models = make(map[string]*model.Model)
	wf.Outputs = make(map[string]string)

	if data.Task.Name != "" {
		wf.Task = data.Task.Name
	}

	// Validate the models
	for k, mcfg := range data.Models {
		fmt.Println("Validating model", k)

		m, err := model.GetModel(mcfg.Identifier)
		if err != nil {
			return nil, err
		}

		// Validate the model options
		err = m.ValidateAndSetOptions(mcfg.Options, cfg)
		if err != nil {
			return nil, err
		}
		wf.Models[k] = &m
	}

	// Validate and build the tools
	for k, tcfg := range data.Tools {
		fmt.Println("Validating tool", k)

		t, err := tool.GetTool(tcfg.Identifier)
		if err != nil {
			return nil, err
		}

		// Validate the model options
		err = t.ValidateAndSetOptions(tcfg.Options, cfg)

		if err != nil {
			return nil, err
		}
		wf.Tools[k] = &t
	}
	if len(data.Task.Steps) == 0 {
		return nil, fmt.Errorf("no steps in the workflow")
	}

	// Create a root node for the workflow
	wf.RootNode = data.Task.Steps[0].ID

	var prev string
	for _, step := range data.Task.Steps {
		// Create the nodes
		node := &WorkflowNode{}
		node.ID = step.ID
		if step.Tool != "" {
			node.Tool = wf.Tools[step.Tool]
			err := (*node.Tool).ValidateParameters(step.Parameters)
			if err != nil {
				return nil, err
			}
			node.Parameters = step.Parameters
		}
		if step.Model != "" {
			node.Model = wf.Models[step.Model]
			err := (*node.Model).ValidateParameters(step.Parameters)
			if err != nil {
				return nil, err
			}
			node.Parameters = step.Parameters
		}
		node.StepName = step.Name
		if node.Tool == nil && node.Model == nil {
			return nil, fmt.Errorf("step %s has no tool or model", step.Name)
		}

		var id string = ""
		if step.ID != "" {
			id = step.ID
			if wf.Graph[id] != nil {
				return nil, fmt.Errorf("duplicate step ID %s", id)
			}
		} else {
			// Generate a unique ID for the step
			for id == "" || wf.Graph[id] == nil {
				id = util.GenerateRandomID(5)
			}
		}
		if prev != "" {
			wf.Graph[prev].Next = id
		}
		prev = id
		wf.Graph[step.ID] = node
	}

	return wf, nil
}
