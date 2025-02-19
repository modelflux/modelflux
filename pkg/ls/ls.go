package ls

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func List() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := path.Join(home, ".modelflux", "workflows")
	repos, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	// Read each direcotry in the workflows directory
	// and list the files in each directory
	for _, repo := range repos {
		if repo.IsDir() {
			workflowDir := path.Join(dir, repo.Name())
			workflowFiles, err := os.ReadDir(workflowDir)
			if err != nil {
				return err
			}
			for _, wf := range workflowFiles {
				// Remove the .yaml extension
				// and print the workflow name {repo}/{workflow}
				wfName, _ := strings.CutSuffix(wf.Name(), ".yaml")
				name := fmt.Sprintf("%s/%s", repo.Name(), wfName)
				fmt.Println(name)
			}
		}
	}
	return nil
}
