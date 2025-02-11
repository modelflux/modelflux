package fileio

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelflux/cli/pkg/util"
)

type fileIOInputs struct {
	Operation string `yaml:"operation"`
	Filepath  string `yaml:"filepath"`
	Content   string `yaml:"content"`
}

type FileIO struct{}

func (t *FileIO) Validate(uParams map[string]interface{}) error {
	params, err := util.BuildStruct[fileIOInputs](uParams)
	if err != nil {
		return err
	}
	operation := params.Operation
	filepath := params.Filepath
	if operation != "read" && operation != "write" {
		return fmt.Errorf("invalid operation: %s", operation)
	}
	if filepath == "" {
		return fmt.Errorf("missing filepath")
	}
	return nil
}

func (t *FileIO) Run(params map[string]interface{}) (string, error) {
	p, err := util.BuildStruct[fileIOInputs](params)
	if err != nil {
		return "", err
	}

	operation := p.Operation
	switch operation {
	case "read":
		return readFile(p.Filepath)
	case "write":
		return writeFile(p.Filepath, p.Content)
	default:
		return "", fmt.Errorf("unsupported operation: %s", operation)
	}
}

func readFile(fp string) (string, error) {
	cleanPath, err := filepath.Abs(fp)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}
	f, err := os.Open(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}
	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}
	return builder.String(), nil
}

func writeFile(fp string, content string) (string, error) {
	cleanPath, err := filepath.Abs(fp)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}
	if err = os.MkdirAll(filepath.Dir(cleanPath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	f, err := os.Create(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	if _, err = writer.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write content: %w", err)
	}
	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush content: %w", err)
	}
	return cleanPath, nil
}
