package pipeline

import (
	"fmt"
	"os"
)

func CustomUsage() {
	fmt.Printf(`Usage: %s [options] <path-to-yaml> [var=value...]

Options:
  -show-scripts    Display scripts that will be executed for matching jobs
  -expand-only     Output expanded YAML configuration without simulation
  -version         Show version information
  -h, -help        Show this help message

Examples:
  %s -show-scripts .gitlab-ci.yml CI_COMMIT_BRANCH=master
  %s -expand-only complex-pipeline.yml > expanded.yml
`, os.Args[0], os.Args[0], os.Args[0])
}
