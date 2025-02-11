package main

import (
	"fmt"

	"github.com/modelflux/cli/pkg/model"
	"github.com/spf13/cobra"
)

var Model string

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Send a message to the model. This is just a test command.",
	Run: func(cmd *cobra.Command, args []string) {
		var input = args[0]
		m, err := model.GetModel(Model)
		if err != nil {
			fmt.Printf("error getting model: %v", err)
			return
		}
		err = m.ValidateAndSetOptions(nil, Config)
		if err != nil {
			fmt.Printf("error validating options: %v", err)
			return
		}

		if err := m.Init(); err != nil {
			fmt.Printf("error initializing model: %v", err)
			return
		}

		resp, err := m.Generate(input)
		if err != nil {
			fmt.Printf("error generating response: %v", err)
			return
		}
		fmt.Println(resp)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&Model, "model", "m", "azure-openai", "Model to use (required)")
}
