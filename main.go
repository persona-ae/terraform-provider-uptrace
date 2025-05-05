package main

import (
	"encoding/json"
	"fmt"
	"os"

	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func prettyPrintResult(response any, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	pretty, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(pretty))
}

func main() {
	projectId := "3255"
	token := "OEkftWB6p3JMXu3MVw9LhA"

	u := uptrace.NewUptraceClient(projectId, token)

	fmt.Println("Getting all...")
	monitors, err := u.GetMonitors()
	prettyPrintResult(monitors, err)

	id := "3592"
	fmt.Printf("Getting by id %s...\n", id)
	monitor, err := u.GetMonitorById(id)
	prettyPrintResult(monitor, err)
}
