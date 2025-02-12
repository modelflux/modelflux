package workflow

import (
	"fmt"

	generate "github.com/modelflux/modelflux/pkg/ai"
	"github.com/modelflux/modelflux/pkg/model"
	"github.com/modelflux/modelflux/pkg/tool"
	"github.com/modelflux/modelflux/pkg/util"
	"github.com/spf13/viper"
)

func (wf *Workflow) ValidateAndBuildWorkflow(data *WorkflowSchema, cfg *viper.Viper) error {
	fmt.Println("VALIDATING WORKFLOW:", data.Name)
	wf.name = data.Name
	wf.graph = make(map[string]*WorkflowNode)
	wf.outputs = make(map[string]string)

	// Create a root node for the workflow
	wf.rootNode = data.Steps[0].ID

	var prev string
	var err error
	for _, step := range data.Steps {
		node := &WorkflowNode{}
		node.ID, err = getStepID(step.ID, wf.graph)
		if err != nil {
			return err
		}
		node.StepName = step.Name
		node.Params = step.With
		node.Log = step.Log
		node.Params = step.With
		node.Operation = step.Run

		// Get the corresponding tool and model.
		if step.Uses == "" && step.Model.Provider == "" {
			return fmt.Errorf("step %s has no tool or model", step.Name)
		} else if step.Uses != "" && step.Model.Provider != "" {
			return fmt.Errorf("step %s has both a tool and a model", step.Name)
		} else if step.Uses != "" {
			t, err := tool.GetTool(step.Uses)
			if err != nil {
				return err
			}
			err = t.Validate(node.Params)
			if err != nil {
				return err
			}

			node.Tool = t
		} else {
			m, err := model.GetModel(step.Model.Provider)
			if err != nil {
				return err
			}
			err = m.ValidateAndSetOptions(step.Model.Options, cfg)
			if err != nil {
				return err
			}
			node.Model = m
			if node.Operation == "generate" {
				err := generate.Validate(node.Params)
				if err != nil {
					return err
				}
			}
		}

		if prev != "" {
			wf.graph[prev].Next = node.ID
		}
		prev = node.ID
		wf.graph[node.ID] = node
	}
	return nil
}

// returns a valid step ID given a provided stepID and the current workflow graph.
// If the provided stepID is non-empty, it checks for duplicates. Otherwise, it generates a new unique ID.
func getStepID(stepID string, steps map[string]*WorkflowNode) (string, error) {
	if stepID != "" {
		if steps[stepID] != nil {
			return "", fmt.Errorf("duplicate step ID %s", stepID)
		}
		return stepID, nil
	}

	var id string
	for {
		id = util.GenerateRandomID(5)
		if steps[id] == nil {
			return id, nil
		}
	}
}
