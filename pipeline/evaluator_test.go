package pipeline

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		job      Job
		expected bool
	}{
		{
			name: "simple_always_rule",
			env:  map[string]string{"CI_COMMIT_BRANCH": "master"},
			job: Job{
				Rules: []Rule{
					{If: `$CI_COMMIT_BRANCH == "master"`, When: "always"},
				},
			},
			expected: true,
		},
		{
			name: "regex_match",
			env:  map[string]string{"CI_COMMIT_TAG": "v1.2.3"},
			job: Job{
				Rules: []Rule{
					{If: `$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/`, When: "manual"},
				},
			},
			expected: false, // Because when: manual
		},
		{
			name: "rule_order_precedence",
			env:  map[string]string{"CI_COMMIT_BRANCH": "master"},
			job: Job{
				Rules: []Rule{
					{If: `$CI_COMMIT_BRANCH == "master"`, When: "never"},
					{If: `$CI_PIPELINE_SOURCE == "schedule"`, When: "always"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvaluator(tt.env)
			result, _ := e.shouldRun(tt.job)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConditionEvaluation(t *testing.T) {
	e := NewEvaluator(map[string]string{
		"CI_COMMIT_BRANCH": "feature/new-auth",
		"DEPLOY_ENV":       "production",
		"CI_COMMIT_TAG":    "v2-45",
	})

	tests := []struct {
		condition string
		expected  bool
	}{
		{`$DEPLOY_ENV == "production" && $CI_COMMIT_TAG =~ /^v\d+-\d+$/`, true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Print transformed condition for debugging
			transformed := e.transformCondition(tt.condition)
			t.Logf("Transformed condition: %s => %s", tt.condition, transformed)

			result, err := e.evaluateCondition(tt.condition)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
