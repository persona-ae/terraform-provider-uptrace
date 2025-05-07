package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/persona-ae/terraform-provider-uptrace/internal/provider"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func old_main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/persona-ae/uptrace",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	c := uptrace.NewUptraceClient("3255", "OEkftWB6p3JMXu3MVw9LhA")

	fmt.Println("Getting by id...")

	var resp uptrace.GetMonitorByIdResponse
	err := c.GetMonitorById(context.Background(), "3592", &resp)
	if err != nil {
		panic(err.Error())
	}
	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(out))

	fmt.Println("Creating monitor...")

	monitor := uptrace.Monitor{
		Name: "tts_p90",
		Type: "metric",
		Params: uptrace.Params{
			Metrics: []uptrace.Metric{
				{
					Name:  "uptrace_tracing_spans",
					Alias: "$spans",
				},
			},
			Query:           "p90($spans) as p90 | where _name = 'stt:finalize'",
			Column:          "",
			MinAllowedValue: 0,
			MaxAllowedValue: 10000,
		},
	}

	var response uptrace.MonitorIdResponse
	err = c.CreateMonitor(context.Background(), monitor, &response)
	if err != nil {
		panic(err.Error())
	}

	out, err = json.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(out))
}
