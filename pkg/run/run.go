package run

import (
	"github.com/modelflux/cli/pkg/workflow"
	"github.com/spf13/viper"
)

type RunFlags struct {
	Local bool
}

func Run(workflowName string, cfg *viper.Viper, flags *RunFlags) {
	schema, err := workflow.LoadSchema(workflowName, flags.Local)
	if err != nil {
		panic(err)
	}

	w := &workflow.Workflow{}
	err = w.ValidateAndBuildWorkflow(schema, cfg)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Initializing workflow")
	err = w.Init()
	if err != nil {
		panic(err)
	}

	err = w.Run()
	if err != nil {
		panic(err)
	}
}
