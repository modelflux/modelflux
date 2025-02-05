package tool

// Tool is the interface all tools must implement.
type Tool interface {
	ValidateOptions(uOptions map[string]interface{}) (interface{}, error)
	ValidateParameters(uParams map[string]interface{}) (interface{}, error)
	Run(cfg interface{}) (interface{}, error)
}
