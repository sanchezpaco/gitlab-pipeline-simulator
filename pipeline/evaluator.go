package pipeline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

type Evaluator struct {
	env map[string]string
}

func NewEvaluator(env map[string]string) *Evaluator {
	return &Evaluator{env: env}
}

func (e *Evaluator) EvaluatePipeline(jobs []Job) map[string][]JobResult {
	results := make(map[string][]JobResult)

	for _, job := range jobs {
		shouldRun, matchedRule := e.shouldRun(job)
		if shouldRun {
			results[job.Stage] = append(results[job.Stage], JobResult{
				Name:             job.Name,
				Stage:            job.Stage,
				MatchedRule:      matchedRule,
				MatchedCondition: matchedRule.If,
				Scripts:          job.Scripts,
			})
		}
	}

	return results
}

func (e *Evaluator) shouldRun(job Job) (bool, Rule) {
	if len(job.Rules) == 0 {
		return true, Rule{}
	}

	for _, rule := range job.Rules {
		result, err := e.evaluateCondition(rule.If)
		if err != nil || !result {
			continue
		}

		switch rule.When {
		case "always":
			return true, rule
		case "never":
			return false, rule
		case "manual":
			return false, rule
		default:
			return true, rule
		}
	}

	return false, Rule{}
}

func (e *Evaluator) evaluateCondition(condition string) (bool, error) {
	if condition == "" {
		return true, nil
	}

	transformed := e.transformCondition(condition)
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(transformed, e.getFunctions())
	if err != nil {
		return false, fmt.Errorf("failed to parse condition: %w", err)
	}

	parameters := make(map[string]interface{})
	for k, v := range e.env {
		parameters[k] = v
	}

	result, err := expression.Evaluate(parameters)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate condition: %w", err)
	}

	return isTruthy(result), nil
}

func (e *Evaluator) transformCondition(condition string) string {
	condition = e.replaceRegexOperators(condition)
	condition = e.substituteVariables(condition)
	return condition
}

func (e *Evaluator) replaceRegexOperators(condition string) string {
	re := regexp.MustCompile(`(\$\w+)\s*=~\s*/((?:\\/|[^/])+)/`)
	return re.ReplaceAllStringFunc(condition, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		varName := parts[1]
		// Escape backslashes in the pattern
		pattern := strings.ReplaceAll(parts[2], `\`, `\\`)
		pattern = strings.ReplaceAll(pattern, `\/`, `/`)
		return fmt.Sprintf("matches(%s, \"%s\")", varName, pattern)
	})
}

func (e *Evaluator) substituteVariables(condition string) string {
	re := regexp.MustCompile(`\$(\w+)`)
	return re.ReplaceAllStringFunc(condition, func(match string) string {
		varName := strings.TrimPrefix(match, "$")
		value, exists := e.env[varName]
		if !exists {
			return `""`
		}
		return fmt.Sprintf("%q", value)
	})
}

func (e *Evaluator) getFunctions() map[string]govaluate.ExpressionFunction {
	return map[string]govaluate.ExpressionFunction{
		"matches": func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return false, fmt.Errorf("matches requires exactly 2 arguments")
			}
			str := fmt.Sprintf("%v", args[0])
			pattern := fmt.Sprintf("%v", args[1])
			matched, err := regexp.MatchString(pattern, str)
			if err != nil {
				return false, fmt.Errorf("invalid regex pattern: %w", err)
			}
			return matched, nil
		},
	}
}

func isTruthy(i interface{}) bool {
	if i == nil {
		return false
	}

	switch v := i.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case int, int8, int16, int32, int64:
		return v != 0
	case uint, uint8, uint16, uint32, uint64:
		return v != 0
	case float32, float64:
		return v != 0.0
	default:
		return true
	}
}
