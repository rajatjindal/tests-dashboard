package gojson

import (
	"regexp"
	"strings"
	"time"

	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

type Event struct {
	Test    string
	Package string
	Action  string
	Time    time.Time
	Elapsed float64
	Output  string
}

type eventProcessor struct {
	suites map[string]*types.Suite
}

func (p *eventProcessor) processEvent(runId string, event Event) {
	suiteName := findSuiteName(event.Test)
	suite, ok := p.suites[suiteName]
	if !ok {
		suite = &types.Suite{
			AllTests:  map[string]*types.Test{},
			TestsTree: []*types.Test{},
		}
	}

	processEvent(suite, runId, suiteName, event)

	p.suites[suiteName] = suite
}

func processEvent(p *types.Suite, runId string, suiteName string, event Event) {
	switch event.Action {
	case "run":
		//suite is starting
		if event.Test == suiteName {
			return
		}

		test := findTest(p, runId, event)
		test.SuiteName = suiteName
		primaryKey := event.Package + "|" + event.Test
		p.AllTests[primaryKey] = test
		p.TestsTree = append(p.TestsTree, test)
	case "output":
		test := findTest(p, runId, event)
		test.Logs += event.Output
	case "pass":
		test := findTest(p, runId, event)
		test.Result = "ok"
		test.Duration = event.Elapsed
	case "fail":
		test := findTest(p, runId, event)
		test.Result = "failed"
		test.Duration = event.Elapsed
	case "skip":
		test := findTest(p, runId, event)
		test.Result = "ignored"
		test.Duration = event.Elapsed
	}
}

func findTest(p *types.Suite, runId string, event Event) *types.Test {
	primaryKey := event.Package + "|" + event.Test
	if test, ok := p.AllTests[primaryKey]; ok {
		return test
	}

	return &types.Test{
		Name:  getHumanReadableName(event.Test),
		RunId: runId,
	}
}

func getHumanReadableName(name string) string {
	if !strings.Contains(name, "/") {
		re := regexp.MustCompile(`[A-Z][^A-Z]*`)
		tokens := re.FindAllString(name, -1)
		return strings.Join(tokens, " ")
	}

	if len(strings.Split(name, "/")) == 2 {
		namex := strings.Split(name, "/")[1]
		return strings.Join(strings.Split(namex, "_"), " ")
	}

	if len(strings.Split(name, "/")) > 2 {
		namex := strings.Split(name, "/")[2]
		return strings.Join(strings.Split(namex, "_"), " ")
	}

	return name
}

func findSuiteName(name string) string {
	if !(strings.Contains(name, "/")) {
		return name
	}

	parts := strings.Split(name, "/")
	return parts[0]
}
