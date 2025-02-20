package run

import (
	"github.com/modelflux/modelflux/pkg/workflow"
	"github.com/spf13/viper"
)

type RunOptions struct {
	Local bool
}

func Run(workflowName string, cfg *viper.Viper, opts *RunOptions) {
	schema, err := workflow.LoadSchema(workflowName, opts.Local)
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
