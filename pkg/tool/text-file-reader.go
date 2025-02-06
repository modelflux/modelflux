package tool

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelflux/cli/pkg/util"
)

type textFileReaderParameters struct {
	Filepath string `yaml:"filepath"`
}
type textFileReaderOptions struct {
}

type TextFileReaderTool struct {
	params  textFileReaderParameters
	options textFileReaderOptions
}

func (t *TextFileReaderTool) ValidateAndSetOptions(uOptions map[string]interface{}) error {
	t.options = textFileReaderOptions{}
	return nil
}

func (t *TextFileReaderTool) ValidateParameters(uParams map[string]interface{}) (interface{}, error) {
	params, err := util.CreateStruct[textFileReaderParameters](uParams)

	if err != nil {
		return nil, err
	}

	return params, nil
}

func (t *TextFileReaderTool) SetParameters(params interface{}) error {
	p := params.(textFileReaderParameters)

	// Additional checks.
	if p.Filepath == "" {
		return fmt.Errorf("missing filepath in text-file-writer parameters")
	}

	t.params = p

	return nil
}

// Run reads the file specified in the parameters and returns its content.
func (t *TextFileReaderTool) Run() (string, error) {
	// Convert the file path to an absolute path.
	cleanPath, err := filepath.Abs(t.params.Filepath)
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
