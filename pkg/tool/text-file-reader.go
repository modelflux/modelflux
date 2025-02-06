package tool

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelflux/cli/pkg/util"
)

type TextFileReaderParams struct {
	Filepath string `yaml:"filepath"`
}

type TextFileReaderTool struct{}

func (t *TextFileReaderTool) ValidateOptions(uOptions map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (t *TextFileReaderTool) ValidateParams(uParams map[string]interface{}) (TextFileReaderParams, error) {
	var params TextFileReaderParams

	// Create a struct from the map using the util package.
	params, err := util.CreateStruct[TextFileReaderParams](uParams)

	if err != nil {
		return TextFileReaderParams{}, err
	}

	// Additional checks
	if params.Filepath == "" {
		return TextFileReaderParams{}, fmt.Errorf("missing filepath in text-file-reader parameters")
	}

	return params, nil

}

// Run reads the file specified in the parameters and returns its content.
func (t *TextFileReaderTool) Run(params TextFileReaderParams) (string, error) {
	// Convert the file path to an absolute path.
	cleanPath, err := filepath.Abs(params.Filepath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	f, err := os.Open(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var content string
	scanner := bufio.NewScanner(f)
	// Read the file line by line. This is done to avoid reading the entire file into memory.
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}

	return string(content), nil
}
