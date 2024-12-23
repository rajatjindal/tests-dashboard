package types

import (
	"strings"
	"time"
)

type Test struct {
	RunId       string          `json:"runId" db:"run_id"`
	SuiteId     string          `json:"suiteId" db:"suite_id"`
	SuiteName   string          `json:"suiteName" db:"suite_name"`
	Name        string          `json:"name" db:"name"`
	Result      string          `json:"result" db:"result"`
	Duration    float64         `json:"duration" db:"duration"`
	Logs        string          `json:"logs" db:"logs"`
	LogsBuilder strings.Builder `json:"-" db:"-"`
	CreatedAt   string          `json:"createdAt" db:"created_at"`
}

type SuiteSummary struct {
	RunId     string  `json:"runId" db:"run_id"`
	SuiteId   string  `json:"suiteId" db:"suite_id"`
	SuiteName string  `json:"suiteName" db:"suite_name"`
	Result    string  `json:"result" db:"result"`
	Passed    int64   `json:"passed" db:"passed"`
	Failed    int64   `json:"failed" db:"failed"`
	Ignored   int64   `json:"ignored" db:"ignored"`
	Duration  float64 `json:"duration" db:"duration"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
}

type Summary struct {
	RunId     string  `json:"runId" db:"run_id"`
	Result    string  `json:"result" db:"result"`
	Passed    int64   `json:"passed" db:"passed"`
	Failed    int64   `json:"failed" db:"failed"`
	Ignored   int64   `json:"ignored" db:"ignored"`
	Duration  float64 `json:"duration" db:"duration"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
}

type Metadata struct {
	RunId     string            `json:"runId" db:"run_id"`
	CommitSha string            `json:"commitSha" db:"commit_sha"`
	JobName   string            `json:"jobName" db:"job_name"`
	Repo      string            `json:"repo" db:"repo"`
	Branch    string            `json:"branch" db:"branch"`
	Format    string            `json:"format" db:"format"`
	Link      string            `json:"link" db:"link"`
	RawTags   string            `json:"-" db:"tags"`
	Tags      map[string]string `json:"tags" db:"-"`
	CreatedAt string            `json:"createdAt" db:"created_at"`
}

type Event struct {
	Type        string  `json:"type"`
	Event       string  `json:"event"`
	Name        string  `json:"name"`
	TestCount   int64   `json:"test_count"`
	Stdout      string  `json:"stdout"`
	Passed      int64   `json:"passed"`
	Failed      int64   `json:"failed"`
	Ignored     int64   `json:"ignored"`
	Measured    int64   `json:"measured"`
	FilteredOut int64   `json:"filtered_out"`
	ExecTime    float64 `json:"exec_time"`
}

type Suite struct {
	RunId     string  `json:"runId" db:"run_id"`
	SuiteId   string  `json:"suiteId" db:"suite_id"`
	SuiteName string  `json:"suiteName" db:"suite_name"`
	Result    string  `json:"result" db:"result"`
	Passed    int64   `json:"passed" db:"passed"`
	Failed    int64   `json:"failed" db:"failed"`
	Ignored   int64   `json:"ignored" db:"ignored"`
	Duration  float64 `json:"duration" db:"duration"`
	CreatedAt string  `json:"createdAt" db:"created_at"`

	StartTime time.Time `json:"startTime" db:"-"`
	EndTime   time.Time `json:"endTime" db:"-"`

	AllTests  map[string]*Test `json:"-" db:"-"`
	TestsTree []*Test          `json:"-" db:"-"`
}

type CommonFilter struct {
	Repo      string            `json:"repo"`
	Branch    string            `json:"branch"`
	CommitSha string            `json:"commitSha"`
	Tags      map[string]string `json:"tags"`
	Page      int               `json:"page"`
	PerPage   int               `json:"perPage"`
}
type SuiteTrendsFilter struct {
	SuiteName string `json:"suiteName"`
	CommonFilter
}

type SuiteTimeTrendEntry struct {
	RunId     string  `json:"runId" db:"run_id"`
	CommitSha string  `json:"commitSha" db:"commit_sha"`
	JobName   string  `json:"jobName" db:"job_name"`
	SuiteName string  `json:"suiteName" db:"suite_name"`
	Duration  float64 `json:"duration" db:"duration"`
}

type Dataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	Stack           string    `json:"stack"`
	BackgroundColor string    `json:"backgroundColor"`
	Hidden          bool      `json:"hidden"`
}
type TimeTrendsData struct {
	Labels   [][]string `json:"labels"`
	Ids      []string   `json:"ids"`
	Datasets []Dataset  `json:"datasets"`
}
