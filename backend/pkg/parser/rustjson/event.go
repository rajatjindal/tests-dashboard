package rustjson

import "github.com/rajatjindal/tests-dashboard/backend/pkg/types"

type eventProcessor struct {
	Summary *types.Summary
	Suites  []types.Suite

	suite *types.Suite
}

func (p *eventProcessor) processEvent(runId string, event types.Event) {
	if event.Type == "suite" && event.Event == "started" {
		p.suite = &types.Suite{
			AllTests: map[string]*types.Test{},
		}
		return
	}

	if event.Type == "suite" && event.Event != "started" {
		// if any suite failed, summary should be failed
		result := p.Summary.Result
		if event.Event != "" && event.Event != "ok" {
			result = event.Event
		}

		p.Summary = &types.Summary{
			RunId:     p.Summary.RunId,
			Result:    result,
			Passed:    p.Summary.Passed + event.Passed,
			Failed:    p.Summary.Failed + event.Failed,
			Ignored:   p.Summary.Ignored + event.Ignored,
			Duration:  p.Summary.Duration + event.ExecTime,
			CreatedAt: p.Summary.CreatedAt,
		}

		//append if there were any tests
		if len(p.suite.AllTests) > 0 {
			p.Suites = append(p.Suites, *p.suite)
		}

		//reset
		p.suite = &types.Suite{
			AllTests: map[string]*types.Test{},
		}

		return
	}

	if event.Type != "test" {
		return
	}

	switch event.Event {
	case "started":
		id := event.Name
		current := &types.Test{Name: event.Name, RunId: runId}
		p.suite.TestsTree = append(p.suite.TestsTree, current)
		p.suite.AllTests[id] = current
	case "output":
		test := p.findTest(runId, event)
		test.Logs = event.Stdout
	case "ok":
		test := p.findTest(runId, event)
		test.Result = "ok"
		test.Duration = event.ExecTime
		test.Logs = event.Stdout
	case "failed":
		test := p.findTest(runId, event)
		test.Result = "failed"
		test.Duration = event.ExecTime
		test.Logs = event.Stdout
	case "ignored":
		test := p.findTest(runId, event)
		test.Result = "ignored"
		test.Duration = event.ExecTime
		test.Logs = event.Stdout
	}
}

func (p *eventProcessor) findTest(runId string, event types.Event) *types.Test {
	primaryKey := event.Name
	if test, ok := p.suite.AllTests[primaryKey]; ok {
		return test
	}

	return &types.Test{
		Name:  getHumanReadableName(event.Name),
		RunId: runId,
	}
}

func getHumanReadableName(name string) string {
	return name

	// if !strings.Contains(name, "/") {
	// 	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	// 	tokens := re.FindAllString(name, -1)
	// 	return strings.Join(tokens, " ")
	// }

	// if len(strings.Split(name, "/")) == 2 {
	// 	namex := strings.Split(name, "/")[1]
	// 	return strings.Join(strings.Split(namex, "_"), " ")
	// }

	// if len(strings.Split(name, "/")) > 2 {
	// 	namex := strings.Split(name, "/")[2]
	// 	return strings.Join(strings.Split(namex, "_"), " ")
	// }

	// return name
}
