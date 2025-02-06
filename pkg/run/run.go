package run

import (
	"github.com/modelflux/cli/pkg/workflow"
	"github.com/spf13/viper"
)

func Run(workflowName string, cfg *viper.Viper) {
	wb := workflow.InitBuilder(cfg)
	w, err := wb.Load(workflowName).ValidateAndBuild()
	if err != nil {
		panic(err)
	}
	w.Init()
	w.Run()
}
