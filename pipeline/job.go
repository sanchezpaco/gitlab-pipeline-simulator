package pipeline

type Job struct {
	Name    string
	Stage   string
	Rules   []Rule
	Scripts []string
}

type Rule struct {
	If   string
	When string
}

type JobResult struct {
	Name             string
	Stage            string
	MatchedRule      Rule
	MatchedCondition string
	Scripts          []string
}
