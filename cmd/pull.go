package main

import (
	"fmt"
	"strings"

	"github.com/modelflux/modelflux/pkg/pull"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a workflow from the server.",
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		splits := strings.Split(s, "/")
		if len(splits) != 2 {
			fmt.Println("Invalid format, input must be in repo/workflow:tag format")
			return
		}
		repo := splits[0]
		splits = strings.Split(splits[1], ":")
		if len(splits) != 2 {
			fmt.Println("Invalid format, input must be in repo/workflow:tag format")
			return
		}
		workflow := splits[0]
		tag := splits[1]
		if err := pull.Pull(repo, workflow, tag, Config); err != nil {
			fmt.Println("Failed to pull workflow: ", err)
			return
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
