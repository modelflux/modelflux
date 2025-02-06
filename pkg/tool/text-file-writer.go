package tool

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelflux/cli/pkg/util"
)

type textFileWriterParameters struct {
	Filepath string `yaml:"filepath"`
	Content  string `yaml:"content"`
}

type textFileWriterOptions struct {
}

type TextFileWriterTool struct {
	params  textFileWriterParameters
	options textFileWriterOptions
}

func (t *TextFileWriterTool) ValidateAndSetOptions(uOptions map[string]interface{}) error {
	t.options = textFileWriterOptions{}
	return nil
}

func (t *TextFileWriterTool) ValidateAndSetParams(uParams map[string]interface{}) error {
	var params textFileWriterParameters

	// CreateStruct
	params, err := util.CreateStruct[textFileWriterParameters](uParams)

	if err != nil {
		return err
	}

	// Additional checks.
	if params.Filepath == "" {
		return fmt.Errorf("missing filepath in text-file-writer parameters")
	}

	t.params = params

	return nil
}

// Run writes the content to the file specified in the parameters.
func (t *TextFileWriterTool) Run() (string, error) {
	// Convert the file path to an absolute path.
	cleanPath, err := filepath.Abs(t.params.Filepath)

	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
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
	_, err = writer.WriteString(t.params.Content)
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return cleanPath, nil
}
