package fetch

import (
	"fmt"
	"io"
	"net/http"

	"github.com/modelflux/cli/pkg/util"
)

type FetchInputs struct {
	URL string `yaml:"url"`
}

type Fetch struct{}

func (f *Fetch) Validate(params map[string]interface{}) error {
	input, err := util.BuildStruct[FetchInputs](params)
	if err != nil {
		return err
	}
	if input.URL == "" {
		return fmt.Errorf("missing url")
	}
	return nil
}

func (f *Fetch) Run(params map[string]interface{}) (string, error) {
	input, err := util.BuildStruct[FetchInputs](params)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(input.URL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(body), nil
}
