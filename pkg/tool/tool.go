package tool

import (
	"fmt"
)

type ToolConfiguration struct {
	Identifier  string                 `yaml:"identifier"`
	ToolOptions map[string]interface{} `yaml:"options,omitempty"`
}

// Tool is the interface all tools must implement.
type Tool interface {
	// Options can be set at build time and are currently static.
	ValidateAndSetOptions(uOptions map[string]interface{}) error
	// Parameters can change between usage, they can be validated with a placeholder for dynamic data at build time
	// and set at runtime.
	ValidateParameters(uParams map[string]interface{}) (interface{}, error)
	SetParameters(params interface{}) error
	Run() (string, error)
}

func ValidateAndBuild(tdata ToolConfiguration) (Tool, error) {
	var t Tool
	switch tdata.Identifier {
	case "text-file-reader":
		t = &TextFileReaderTool{}
	case "text-file-writer":
		t = &TextFileWriterTool{}
	default:
		return nil, fmt.Errorf("model %s not found", tdata.Identifier)
	}

	if err := t.ValidateAndSetOptions(tdata.ToolOptions); err != nil {
		return nil, fmt.Errorf("error validating model options: %v", err)
	}

	return t, nil
}
