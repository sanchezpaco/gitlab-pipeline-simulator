package pipeline

import (
	"fmt"
	"os"
)

type Simulator struct {
	evaluator   *Evaluator
	showScripts bool
}

func NewSimulator(env map[string]string, showScripts bool) *Simulator {
	return &Simulator{
		evaluator:   NewEvaluator(env),
		showScripts: showScripts,
	}
}

func (s *Simulator) Run(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	root := expandYAML(data)
	jobs, stageOrder, err := extractJobsAndStages(root)
	if err != nil {
		return fmt.Errorf("failed to extract jobs and stages: %w", err)
	}

	results := s.evaluator.EvaluatePipeline(jobs)
	PrintResults(results, stageOrder, s.showScripts)

	return nil
}
