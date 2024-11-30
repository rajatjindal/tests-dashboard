package gojson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

func Ingest(runId string, data []byte) (*types.Summary, []types.Suite, error) {
	p := &eventProcessor{
		suites: map[string]*types.Suite{},
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	index := 0
	for scanner.Scan() {
		index++
		var event Event

		text := scanner.Text()

		if strings.TrimSpace(text) == "" {
			continue
		}

		err := json.Unmarshal([]byte(text), &event)
		if err != nil {
			fmt.Println(index, " ", text)
			fmt.Printf("ignore error index: %d, txt: '%s', err: %v\n", index, text, err)
			continue
		}

		p.processEvent(runId, event)
	}

	summary := &types.Summary{
		RunId:    runId,
		Result:   "pass", //assume pass to start with
		Passed:   0,
		Failed:   0,
		Ignored:  0,
		Duration: 0,
	}

	allSuites := []types.Suite{}
	for _, s := range p.suites {
		if summary.CreatedAt == "" {
			// start time of first suite
			// while this is not accurate, but should be good enough
			summary.CreatedAt = s.CreatedAt
		}

		allSuites = append(allSuites, *s)
		for _, t := range s.TestsTree {
			switch t.Result {
			case "ok":
				summary.Passed = summary.Passed + 1
			case "failed":
				summary.Failed = summary.Failed + 1
			case "ignored":
				summary.Ignored = summary.Ignored + 1
			}
		}
	}

	return summary, allSuites, nil
}
