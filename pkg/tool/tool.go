package tool

type ToolConfiguration struct {
	Identifier  string                 `yaml:"identifier"`
	ToolOptions map[string]interface{} `yaml:"options,omitempty"`
}

// Tool is the interface all tools must implement.
type Tool interface {
	ValidateAndSetOptions(uOptions map[string]interface{}) error
	ValidateAndSetParameters(uParams map[string]interface{}) error
	Run() (string, error)
}
