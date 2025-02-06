package workflow

import (
	"fmt"
	"regexp"

	"github.com/modelflux/cli/pkg/load"
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

// WorkflowBuilder provides a chainable API to build a Workflow.
type WorkflowBuilder struct {
	data *load.WorkflowSchema
	cfg  *viper.Viper
	err  error
}

func InitBuilder(cfg *viper.Viper) *WorkflowBuilder {
	return &WorkflowBuilder{cfg: cfg}
}

// Parse parses the YAML into the workflow structure.
func (b *WorkflowBuilder) Load(workflowName string) *WorkflowBuilder {
	if b.err != nil {
		return b
	}
	b.data = load.Load(workflowName)
	return b
}

func (b *WorkflowBuilder) ValidateAndBuild() (*Workflow, error) {
	fmt.Println()
	if b.err != nil {
		return nil, b.err
	}

	if b.data == nil {
		return nil, fmt.Errorf("workflow data is empty")
	}

	// Initialize the workflow
	wf := &Workflow{}
	wf.graph = make(map[string]*WorkflowNode)
	wf.models = make(map[string]*model.Model)
	wf.tools = make(map[string]*tool.Tool)

	if b.data.Task.Name != "" {
		wf.task = b.data.Task.Name
	}

	for k, m := range b.data.Models {
		fmt.Println("Validating model", k)
		// Build the model
		built_m, err := model.ValidateAndBuild(m, b.cfg)
		if err != nil {
			b.err = err
			return nil, err
		}

		wf.models[k] = &built_m
	}

	for k, t := range b.data.Tools {
		fmt.Println("Validating tool", k)
		// Build the tool
		built_tool, err := tool.ValidateAndBuild(t)
		if err != nil {
			b.err = err
			return nil, err
		}
		wf.tools[k] = &built_tool
	}
	if len(b.data.Task.Steps) == 0 {
		b.err = fmt.Errorf("no steps in workflow")
		return nil, b.err
	}

	// Create a root node for the workflow
	wf.rootNode = b.data.Task.Steps[0].ID

	var prev string
	for _, step := range b.data.Task.Steps {
		// Create the nodes
		node := &WorkflowNode{}
		node.id = step.ID
		node.tool = wf.tools[step.Tool]
		node.model = wf.models[step.Model]
		node.stepName = step.Name
		if node.tool == nil && node.model == nil {
			b.err = fmt.Errorf("step %s has no tool or model", step.Name)
			return nil, b.err
		}
		if node.tool != nil {
			t := *node.tool
			p, err := t.ValidateParameters(step.Parameters)
			if err != nil {
				b.err = err
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
