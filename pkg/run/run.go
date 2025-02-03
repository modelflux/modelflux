package run

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/orelbn/tbd/pkg/load"
)

func Run(workflowName string) {
	// Get the workflow file
	fmt.Println("run called")
	
	workflow := load.Load(workflowName) // or your workflow name
    if workflow == nil {
        log.Fatal("Workflow loading failed.")
    }

    pretty, err := json.MarshalIndent(workflow, "", "  ")
    if err != nil {
        log.Fatal("Error pretty printing:", err)
    }

    fmt.Println(string(pretty))
}
