package main

import (
	"github.com/modelflux/cli/pkg/ls"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available workflows",
	Run: func(cmd *cobra.Command, args []string) {
		ls.List()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
