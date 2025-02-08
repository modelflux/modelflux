package generate

import (
	"fmt"
	"strings"

	"github.com/modelflux/cli/pkg/model"
	"github.com/modelflux/cli/pkg/util"
)

type GenerateParameters struct {
	Prompt   string            `yaml:"prompt"`
	Template string            `yaml:"template"`
	Vars     map[string]string `yaml:"vars"`
}

func Validate(uParams map[string]any) error {
	p, err := util.BuildStruct[GenerateParameters](uParams)
	if err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}
	if p.Prompt == "" && p.Template == "" {
		return fmt.Errorf("prompt or template required")
	}
	return nil
}

func Run(params map[string]any, m model.Model) (string, error) {
	p, err := util.BuildStruct[GenerateParameters](params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	var input string
	if p.Prompt != "" {
		input = p.Prompt
	} else if p.Template != "" {
		input = ReplacePlaceholders(p.Template, p.Vars)
	}
	return m.Generate(input)

}

func ReplacePlaceholders(input string, vars map[string]string) string {
	for k, v := range vars {
		input = strings.ReplaceAll(input, fmt.Sprintf("{%s}", k), v)
	}
	return input
}
