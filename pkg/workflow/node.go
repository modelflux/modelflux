package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	generate "github.com/modelflux/modelflux/pkg/ai"
	"github.com/modelflux/modelflux/pkg/model"
	"github.com/modelflux/modelflux/pkg/tool"
)

// A WorkflowNode represents a step in the workflow.
type WorkflowNode struct {
	StepName  string
	ID        string
	Params    map[string]interface{}
	Operation string
	Tool      tool.Tool
	Model     model.Model
	Next      string // The ID of the next node to run
	Output    string
	Log       bool
}

// replacePlaceholders processes the WorkflowNode's Params by scanning for placeholders
// in the JSON-encoded representation of the Params map and substituting them with values
// from the provided outputs map.
//
// Placeholders must follow the format:
//
//	${{ <identifier>.output }}
//
// The <identifier> component can include alphanumeric characters and dashes.
// For each detected placeholder, the function looks up the corresponding value from outputs.
// The replacement value is JSON-escaped to ensure that any special characters are handled properly.
//
// After all replacements are completed, the modified JSON string is unmarshaled back into the Params map.
// The function returns an error if marshalling or unmarshalling fails.
func (n *WorkflowNode) replacePlaceholders(outputs map[string]string) error {
	// Marshal the entire Params map.
	b, err := json.Marshal(n.Params)
	if err != nil {
		return fmt.Errorf("failed to marshal Params: %v", err)
	}
	asStr := string(b)

	// Regexp that allows spaces in the placeholder name.
	// Example matching: ${{ somestep.output }}
	placeholderRe := regexp.MustCompile(`\${{\s*([a-zA-Z0-9-]+)\.output\s*}}`)

	// Find and replace all placeholders.
	asStr = placeholderRe.ReplaceAllStringFunc(asStr, func(match string) string {
		// Extract the identifier.
		submatches := placeholderRe.FindStringSubmatch(match)
		if len(submatches) != 2 {
			return match
		}
		key := strings.TrimSpace(submatches[1])
		if replacement, ok := outputs[key]; ok {
			// Marshal the replacement to correctly escape special characters.
			escaped, err := json.Marshal(replacement)
			if err != nil {
				// In case of error, use the raw replacement.
				return replacement
			}
			// json.Marshal returns a quoted string, so remove the first and last character.
			escapedStr := string(escaped)
			if len(escapedStr) >= 2 {
				return escapedStr[1 : len(escapedStr)-1]
			}
			return escapedStr
		}
		return match
	})

	// Unmarshal the replaced JSON string back into n.Params.
	var newParams map[string]interface{}
	if err := json.Unmarshal([]byte(asStr), &newParams); err != nil {
		return fmt.Errorf("failed to unmarshal modified Params: %v", err)
	}
	n.Params = newParams
	return nil
}
func (n *WorkflowNode) Run(outputs map[string]string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stdout))
	s.Prefix = "PROCESSING STEP: " + n.StepName + " "
	s.Start()

	fmt.Println("-------------------------")
	fmt.Println("STEP:", n.StepName)
	var output string
	var err error

	// Replace placeholders in the parameters
	err = n.replacePlaceholders(outputs)
	if err != nil {
		s.FinalMSG = "STEP FAILED ❌\n"
		s.Stop()
		return "", err
	}

	if n.Tool != nil {
		output, err = n.Tool.Run(n.Params)
	} else if n.Model != nil {
		if n.Operation == "generate" {
			output, err = generate.Run(n.Params, n.Model)
		}
	}

	if err != nil {
		s.FinalMSG = "STEP FAILED ❌\n"
		s.Stop()
		return "", err
	}

	n.Output = output
	s.FinalMSG = "COMPLETED SUCCESSFULLY ✅\n"
	s.Stop()
	if n.Log {
		fmt.Println()
		fmt.Println(n.Output)
		fmt.Println()
	}
	return n.Next, nil
}
