package main

import (
	"flag"
	"fmt"
	"gitlab-pipeline-simulator/pipeline"
	"log"
	"os"
	"strings"
)

var (
	showScripts  bool
	expandOnly   bool
	printVersion bool
	version      = "dev" // Set during build with ldflags
)

func main() {
	flag.BoolVar(&showScripts, "show-scripts", false, "Display scripts that will be executed for each job")
	flag.BoolVar(&expandOnly, "expand-only", false, "Output expanded YAML configuration without simulation")
	flag.BoolVar(&printVersion, "version", false, "Show version information")
	flag.Usage = pipeline.CustomUsage
	flag.Parse()

	if printVersion {
		fmt.Printf("pipeline-simulator version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filePath := flag.Arg(0)
	env := parseEnv(flag.Args()[1:])

	if expandOnly {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		root := pipeline.ExpandYAML(data)
		pipeline.PrintExpandedYAML(root)
		os.Exit(0)
	}

	simulator := pipeline.NewSimulator(env, showScripts)
	if err := simulator.Run(filePath); err != nil {
		log.Fatalf("Error running pipeline simulation: %v", err)
		log.Fatal(err)
	}
}

func parseEnv(args []string) map[string]string {
	env := make(map[string]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}
