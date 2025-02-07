package pull

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Pull(repo string, workflow string, tag string, cfg *viper.Viper) error {
	fmt.Println("Pulling workflow", repo, workflow, tag)
	registryUrl := cfg.GetString("registryUrl")
	url := registryUrl + "/workflow/pull?repo=" + repo + "&workflow=" + workflow + "&tag=" + tag
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to pull workflow: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull workflow: %s", resp.Status)
	}

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	saveDir := path.Join(home, ".modelflux", "workflows", repo)
	fileName := workflow + ":" + tag + ".yaml"

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path.Join(saveDir, fileName), respBody, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save workflow: %v", err)
	}

	fmt.Printf("Workflow %s:%s saved to %s\n", workflow, tag, path.Join(saveDir, fileName))
	return nil
}
