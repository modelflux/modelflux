package tool

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelflux/cli/pkg/util"
)

type TextFileWriterParameters struct {
	Filepath string `yaml:"filepath"`
	Content  string `yaml:"content"`
}

type TextFileWriterTool struct{}

func (t *TextFileWriterTool) ValidateOptions(uOptions map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (t *TextFileWriterTool) ValidateParams(uParams map[string]interface{}) (TextFileWriterParameters, error) {
	var params TextFileWriterParameters

	// CreateStruct
	params, err := util.CreateStruct[TextFileWriterParameters](uParams)

	if err != nil {
		return TextFileWriterParameters{}, err
	}

	// Additional checks.
	if params.Filepath == "" {
		return TextFileWriterParameters{}, fmt.Errorf("missing filepath in text-file-writer parameters")
	}

	return params, nil
}

// Run writes the content to the file specified in the parameters.
func (t *TextFileWriterTool) Run(params TextFileWriterParameters) (string, error) {

	// Convert the file path to an absolute path.
	cleanPath, err := filepath.Abs(params.Filepath)
	if err != nil {
		return "", err
	}

	// Ensure that the directory exists.
	err = os.MkdirAll(filepath.Dir(cleanPath), os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(cleanPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(params.Content)
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return cleanPath, nil
}
