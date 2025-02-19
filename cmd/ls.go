package main

import (
	"fmt"

	"github.com/modelflux/modelflux/pkg/ls"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available workflows",
	Run: func(cmd *cobra.Command, args []string) {
		err := ls.List()
		if err != nil {
			fmt.Println("Failed to list workflows:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
