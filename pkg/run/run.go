package run

import (
	"github.com/modelflux/cli/pkg/workflow"
	"github.com/spf13/viper"
)

func Run(workflowName string, cfg *viper.Viper) {
	schema, err := workflow.LoadSchema(workflowName)
	if err != nil {
		panic(err)
	}
	w, err := workflow.ValidateAndBuildWorkflow(schema, cfg)
	if err != nil {
		panic(err)
	}

	err = w.Init()
	if err != nil {
		panic(err)
	}

	err = w.Run()
	if err != nil {
		panic(err)
	}
}
