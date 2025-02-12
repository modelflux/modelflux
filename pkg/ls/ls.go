package ls

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func List() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		return
	}
	dir := path.Join(home, ".modelflux", "workflows")
	repos, err := os.ReadDir(dir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read directory: %v\n", err)
		return
	}

	// Read each direcotry in the workflows directory
	// and list the files in each directory
	for _, repo := range repos {
		if repo.IsDir() {
			workflowDir := path.Join(dir, repo.Name())
			workflowFiles, err := os.ReadDir(workflowDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read directory: %v\n", err)
				return
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
}
