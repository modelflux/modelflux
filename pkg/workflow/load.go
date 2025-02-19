package workflow

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/modelflux/modelflux/pkg/config"
	"github.com/modelflux/modelflux/pkg/model"
	"gopkg.in/yaml.v3"
)

type WorkflowSchema struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	ID    string                   `yaml:"id,omitempty"`
	Name  string                   `yaml:"name"`
	Uses  string                   `yaml:"uses,omitempty"`
	Model model.ModelConfiguration `yaml:"model,omitempty"`
	Run   string                   `yaml:"run,omitempty"`  // Run is an operation to be preformed by the model
	With  map[string]interface{}   `yaml:"with,omitempty"` // With is the parameters to be passed to the tool.
	Log   bool                     `yaml:"log,omitempty"`  // Log is a flag wether the output of the tool should be logged to the console.
}

func LoadSchema(workflowName string, local bool) (*WorkflowSchema, error) {
	fmt.Println("LOADING WORKFLOW:", workflowName)

	data, err := readFile(workflowName+".yaml", local)
	if err != nil {
		return nil, err
	}

	var workflow WorkflowSchema
	err = yaml.Unmarshal(data, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func readFile(filename string, local bool) ([]byte, error) {
	var rootDir string
	if local {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		rootDir = cwd
	} else {
		workflowsDir, err := config.GetWorkflowsPath()
		if err != nil {
			return nil, err
		}
		rootDir = workflowsDir
	}
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return nil, err
	}

	defer root.Close()

	file, err := root.Open(filename)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) && strings.Contains(pathErr.Error(), "path escapes from parent") {
			if local {
				return nil, fmt.Errorf("Can not read workflows outside of the working directory: %s", filename)
			} else {
				return nil, fmt.Errorf("Can not read workflows outside of the workflow directory: %s", filename)
			}
		}
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil

}
