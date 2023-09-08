package junit

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/rajatjindal/test-dashboard/backend/pkg/types"
)

func Ingest(runId string, data []byte) (*types.Summary, []types.Suite, error) {
	summary := &types.Summary{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	index := 0
	suites := []Testsuites{}
	for scanner.Scan() {
		index++
		var event Testsuites

		text := scanner.Text()
		err := xml.Unmarshal([]byte(text), &event)
		if err != nil {
			fmt.Println(index, " ", text)
			fmt.Printf("ignore error index: %d, txt: %s, err: %v", index, text, err)
			continue
		}
		suites = append(suites, event)
	}

	suitesx := []types.Suite{}
	for index, suite := range suites {
		summary = &types.Summary{
			RunId:     runId,
			Result:    getResultForSummary(summary.Result, suite),
			Passed:    summary.Passed + getPassed(suite),
			Failed:    summary.Failed + mustInt(suite.Testsuite.Failures),
			Ignored:   summary.Ignored + mustInt(suite.Testsuite.Skipped),
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		tests := []*types.Test{}
		if suite.Testsuite.Tests == "0" {
			continue
		}

		for _, tc := range suite.Testsuite.Testcases {
			tests = append(tests, &types.Test{
				RunId:     runId,
				SuiteId:   fmt.Sprintf("%d", index),
				Name:      fmt.Sprintf("%s::%s\n", tc.Classname, tc.Name),
				Result:    getResult(tc),
				Duration:  mustFloat64(tc.Time),
				Logs:      tc.SystemOut,
				CreatedAt: time.Now().Format(time.RFC3339),
			})
		}

		suitesx = append(suitesx, types.Suite{
			TestsTree: tests,
		})
	}

	return summary, suitesx, nil
}

func getResult(tc TestCase) string {
	if tc.Failure != nil {
		return "failed"
	}

	return "ok"
}

func mustFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func mustInt(s string) int64 {
	f, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return int64(f)
}

func getResultForSummary(result string, suite Testsuites) string {
	if result != "ok" {
		if mustInt(suite.Testsuite.Failures) > 0 {
			return "failed"
		}
	}

	return "ok"
}

func getPassed(suite Testsuites) int64 {
	return mustInt(suite.Testsuite.Tests) - (mustInt(suite.Testsuite.Failures) + mustInt(suite.Testsuite.Skipped))
}
