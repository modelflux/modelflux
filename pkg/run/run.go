package run

import (
	"github.com/modelflux/modelflux/pkg/workflow"
	"github.com/spf13/viper"
)

func Run(workflowName string, cfg *viper.Viper) {
	schema, err := workflow.LoadSchema(workflowName)
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
