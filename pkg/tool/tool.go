package tool

import (
	"fmt"

	"github.com/spf13/viper"
)

type ToolConfiguration struct {
	Identifier string                 `yaml:"identifier"`
	Options    map[string]interface{} `yaml:"options,omitempty"`
}

// Tool is the interface all tools must implement.
type Tool interface {
	// Options can be set at build time and are currently static.
	ValidateAndSetOptions(uOptions map[string]interface{}, cfg *viper.Viper) error
	New() error
	// Parameters can change between usage, they can be validated with a placeholder for dynamic data at build time
	// and set at runtime.
	ValidateParameters(uParams map[string]interface{}) error
	SetParameters(uParams map[string]interface{}) error
	Run() (string, error)
}

func GetTool(t string) (Tool, error) {
	switch t {
	case "text-file-reader":
		return &TextFileReaderTool{}, nil
	case "text-file-writer":
		return &TextFileWriterTool{}, nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", t)
	}
}
