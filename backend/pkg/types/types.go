package types

type Test struct {
	RunId     string  `json:"runId" db:"run_id"`
	SuiteId   string  `json:"suiteId" db:"suite_id"`
	Name      string  `json:"name" db:"name"`
	Result    string  `json:"result" db:"result"`
	Duration  float64 `json:"duration" db:"duration"`
	Logs      string  `json:"logs" db:"logs"`
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
	RunId     string `json:"runId" db:"run_id"`
	Repo      string `json:"repo" db:"repo"`
	Branch    string `json:"branch" db:"branch"`
	Format    string `json:"format" db:"format"`
	Link      string `json:"link" db:"link"`
	Tags      string `json:"tags" db:"tags"`
	CreatedAt string `json:"createdAt" db:"created_at"`
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
	AllTests  map[string]*Test `json:"-" db:"-"`
	TestsTree []*Test
}
