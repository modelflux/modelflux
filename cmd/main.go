package main

import (
	"os"

	"github.com/modelflux/cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	Execute()
}

var Config = viper.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mf",
	Short: "modelflux: Run hundreds of AI workflows right from the command line.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.InitConfig(Config)
}
