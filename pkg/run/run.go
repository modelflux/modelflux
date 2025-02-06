package run

import (
	"github.com/modelflux/cli/pkg/workflow"
	"github.com/spf13/viper"
)

func Run(workflowName string, cfg *viper.Viper) {
	schema, err := workflow.Load(workflowName)
	if err != nil {
		panic(err)
	}
	w, err := workflow.ValidateAndBuild(schema, cfg)
	if err != nil {
		panic(err)
	}

	w.Init()
	w.Run()
}
