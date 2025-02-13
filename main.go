package main

import (
	"flag"
	"gitlab-pipeline-simulator/pipeline"
	"log"
	"strings"
)

var showScripts bool

func main() {
	flag.BoolVar(&showScripts, "show-scripts", false, "Display scripts that will be executed for each job")
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("Usage: ./pipeline-simulator [flags] <path-to-yaml> <var1=value1> <var2=value2> ...")
	}

	env := parseEnv(flag.Args()[1:])
	simulator := pipeline.NewSimulator(env, showScripts)

	filePath := flag.Arg(0)
	if err := simulator.Run(filePath); err != nil {
		log.Fatalf("Error running pipeline simulation: %v", err)
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
