package workflow

import (
	"fmt"
	"regexp"

	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/tool"
	"github.com/modelflux/cli/pkg/util"
	"github.com/spf13/viper"
)

type WorkflowNode struct {
	stepName   string
	id         string
	depends_on []string // Identifiers of dependent nodes
	parameters interface{}
	tool       *tool.Tool
	model      *model.Model
	Output     string // The result of running this node
	next       string
}

// Workflow represents your parsed and built workflow.
type Workflow struct {
	task     string
	graph    map[string]*WorkflowNode
	tools    map[string]*tool.Tool
	models   map[string]*model.Model
	rootNode string
}

func (n *WorkflowNode) Run() (string, error) {
	fmt.Println("Running step:", n.stepName)
	// Do somethings
	n.Output = "output " + n.id
	if n.next != "" {
	} else {
		fmt.Println("Done")
	}
	return n.next, nil
}

func ValidateAndBuild(data *WorkflowSchema, cfg *viper.Viper) (*Workflow, error) {
	fmt.Println()

	// Initialize the workflow
	wf := &Workflow{}
	wf.graph = make(map[string]*WorkflowNode)
	wf.models = make(map[string]*model.Model)
	wf.tools = make(map[string]*tool.Tool)

	if data.Task.Name != "" {
		wf.task = data.Task.Name
	}

	for k, m := range data.Models {
		fmt.Println("Validating model", k)
		// Build the model
		built_m, err := model.ValidateAndBuild(m, cfg)
		if err != nil {
			return nil, err
		}

		wf.models[k] = &built_m
	}

	for k, t := range data.Tools {
		fmt.Println("Validating tool", k)
		// Build the tool
		built_tool, err := tool.ValidateAndBuild(t)
		if err != nil {
			return nil, err
		}
		wf.tools[k] = &built_tool
	}
	if len(data.Task.Steps) == 0 {
		return nil, fmt.Errorf("no steps in the workflow")
	}

	// Create a root node for the workflow
	wf.rootNode = data.Task.Steps[0].ID

	var prev string
	for _, step := range data.Task.Steps {
		// Create the nodes
		node := &WorkflowNode{}
		node.id = step.ID
		node.tool = wf.tools[step.Tool]
		node.model = wf.models[step.Model]
		node.stepName = step.Name
		if node.tool == nil && node.model == nil {
			return nil, fmt.Errorf("step %s has no tool or model", step.Name)
		}
		if node.tool != nil {
			t := *node.tool
			p, err := t.ValidateParameters(step.Parameters)
			if err != nil {
				return nil, err
			}
			node.parameters = p
			// Check if any of the parameters have a placeholder
			// and record it to be replaced at runtime
			for _, value := range step.Parameters {
				println("Checking for placeholders in", value)
				var placeholderRegex = regexp.MustCompile(`\{\{([\w-]+)\.output\}\}`)
				if strValue, ok := value.(string); ok && placeholderRegex.MatchString(strValue) {
					// Extract the id of the step that will provide the output
					// for this placeholder this is the portion before the dot
					// in the placeholder string this is the portion before the .output
					dependency := placeholderRegex.FindStringSubmatch(strValue)[0]
					println("Found dependency", dependency)
					node.depends_on = append(node.depends_on, dependency)
				}
			}
		}
		var id string = ""
		if step.ID != "" {
			id = step.ID
			if wf.graph[id] != nil {
				return nil, fmt.Errorf("duplicate step ID %s", id)
			}
		} else {
			// Generate a unique ID for the step
			for id == "" || wf.graph[id] == nil {
				id = util.GenerateRandomID(5)
			}
		}
		if prev != "" {
			wf.graph[prev].next = id
		}
		prev = id
		wf.graph[step.ID] = node
	}

	return wf, nil
}
