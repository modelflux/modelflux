package tool

import (
	"fmt"

	"github.com/modelflux/modelflux/pkg/fetch"
	"github.com/modelflux/modelflux/pkg/fileio"
	"github.com/modelflux/modelflux/pkg/htmlparser"
)

type ToolConfiguration struct {
	Source  string         `yaml:"source"`
	Options map[string]any `yaml:"options,omitempty"`
}

// Tool is the interface all tools must implement.
type Tool interface { //params are a stringified yaml
	Validate(uParams map[string]interface{}) error
	Run(params map[string]interface{}) (string, error)
}

func GetTool(name string) (Tool, error) {
	switch name {
	case "modelflux/fileio":
		return &fileio.FileIO{}, nil
	case "modelflux/fetch":
		return &fetch.Fetch{}, nil
	case "modelflux/htmlparser":
		return &htmlparser.HTMLParser{}, nil
	default:
		return nil, fmt.Errorf("tool %s not found", name)
	}
}
