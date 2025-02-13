package pipeline

import (
	"fmt"
)

func PrintResults(results map[string][]JobResult, stageOrder []string, showScripts bool) {
	fmt.Println("\n=== Pipeline Simulation Results ===")
	fmt.Println("Jobs that would run based on current conditions:")

	for _, stage := range stageOrder {
		if jobResults, exists := results[stage]; exists {
			fmt.Printf("🚀 Stage: %s\n", stage)
			for _, result := range jobResults {
				fmt.Printf("   ✅ Job: %s\n", result.Name)
				if result.MatchedCondition != "" {
					fmt.Printf("      ├─ Condition: %s\n", colorizeCondition(result.MatchedCondition))
				}
				if showScripts && len(result.Scripts) > 0 {
					fmt.Println("      ├─ Scripts:")
					for _, script := range result.Scripts {
						fmt.Printf("      │  - %s\n", script)
					}
				}
				fmt.Println()
			}
		}
	}
}

func colorizeCondition(cond string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", cond) // Cyan
}
