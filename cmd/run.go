package main

import (
	"github.com/modelflux/modelflux/pkg/run"
	"github.com/spf13/cobra"
)

var Local bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [workflowName]",
	Short: "Runs a workflow file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		flags := &run.RunOptions{
			Local: Local,
		}

		run.Run(name, Config, flags)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add local flag -l for running locally.
	runCmd.Flags().BoolVarP(&Local, "local", "l", false, "Run a local workflow file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
