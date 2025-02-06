package run

import (
	"github.com/modelflux/cli/pkg/workflow"
	"github.com/spf13/viper"
)

func Run(workflowName string, cfg *viper.Viper) {
	w, err := workflow.BuildWorkflow(workflowName, cfg).Load().Validate().Build()
	if err != nil {
		panic(err)
	}
	w.Run()
}
