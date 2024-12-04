package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/fermyon/spin-go-sdk/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

const defaultPerPage = 10

func db() *sqlx.DB {
	conn := sqlite.Open("default")
	return sqlx.NewDb(conn, "sqlite")
}

func FetchAllTests(ctx context.Context, filter types.CommonFilter) ([]*types.Test, error) {
	conn := db()
	if conn == nil {
		return nil, fmt.Errorf("failed to read db")
	}
	defer conn.Close()

	query := "select tests.* from tests, metadata where metadata.run_id = tests.run_id"
	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	rows, err := conn.Queryx(query, params...)
	if err != nil {
		return nil, err
	}

	tests := []*types.Test{}
	for rows.Next() {
		var test types.Test
		err = rows.StructScan(&test)
		if err != nil {
			return nil, err
		}

		tests = append(tests, &test)
	}

	return tests, nil
}

func FetchSuitesForRun(ctx context.Context, runId string) ([]string, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select distinct(suite_id) from tests where run_id = ?;", runId)
	if err != nil {
		return nil, err
	}

	items := []string{}
	for rows.Next() {
		var item string
		err = rows.Scan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func FetchHistoryForLogLine(ctx context.Context, logLine string, filter types.CommonFilter) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	query := "select tests.* from tests, metadata where metadata.run_id = tests.run_id AND lower(tests.logs) LIKE ?"
	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	rows, err := conn.Queryx(query, append([]interface{}{"%" + logLine + "%"}, params...)...)
	if err != nil {
		return nil, err
	}

	tests := []*types.Test{}
	for rows.Next() {
		var test types.Test
		err = rows.StructScan(&test)
		if err != nil {
			return nil, err
		}

		tests = append(tests, &test)
	}

	return tests, nil
}

func FetchHistoryForTestcase(ctx context.Context, testname string, filter types.CommonFilter) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	query := "select tests.* from tests, metadata where metadata.run_id = tests.run_id AND lower(tests.name) like ?"
	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	perPage := filter.PerPage

	//defaults to defaultPerPage if filter.PerPage is not provided
	if perPage == 0 {
		perPage = defaultPerPage
	}

	query += " order by metadata.created_at desc"
	if filter.Page > 0 {
		query += " offset ?"
		params = append(params, filter.Page*perPage)
	}

	if perPage > 0 {
		query += " limit ?"
		params = append(params, perPage)
	}

	fmt.Println(query)
	fmt.Printf("%#v\n", append([]interface{}{"%" + strings.ToLower(testname) + "%"}, params...))
	rows, err := conn.Queryx(query, append([]interface{}{strings.ToLower(testname)}, params...)...)
	if err != nil {
		return nil, err
	}

	tests := []*types.Test{}
	for rows.Next() {
		var test types.Test
		err = rows.StructScan(&test)
		if err != nil {
			return nil, err
		}

		test.Logs = ""
		tests = append(tests, &test)
	}

	return tests, nil
}

func FetchSuiteSummaryForRunId(ctx context.Context, runId string) ([]*types.SuiteSummary, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx(`
		select 
			run_id,
			suite_id,
			suite_name, 
			passed,
			failed,
			ignored,
			duration
		from 
			suite_summary
		where 
			run_id = ?;`, runId)
	if err != nil {
		return nil, err
	}

	suites := []*types.SuiteSummary{}
	for rows.Next() {
		var summary types.SuiteSummary
		err = rows.StructScan(&summary)
		if err != nil {
			return nil, err
		}

		suites = append(suites, &summary)
	}

	return suites, nil
}

// TODO(rajatjindal): fetch failed only on server side to improve perf of this query
func FetchTestsByRunIdAndSuite(ctx context.Context, runId, suiteId string) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from tests where run_id = ? AND suite_id = ?;", runId, suiteId)
	if err != nil {
		return nil, err
	}

	tests := []*types.Test{}
	for rows.Next() {
		var test types.Test
		err = rows.StructScan(&test)
		if err != nil {
			return nil, err
		}

		tests = append(tests, &test)
	}

	return tests, nil
}

func FetchSummaryForRun(ctx context.Context, runId string) (*types.Summary, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from summary where run_id = ? LIMIT 1;", runId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Summary
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		return &item, nil
	}

	return nil, fmt.Errorf("no summary found for run id %s", runId)
}

func FetchMetadataForRun(ctx context.Context, runId string) (*types.Metadata, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from metadata where run_id = ? LIMIT 1;", runId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Metadata
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var tags map[string]string
		err = json.Unmarshal([]byte(item.RawTags), &tags)
		if err != nil {
			tags["ERROR"] = "FAILED TO PARSE TAGS"
		}

		item.Tags = tags
		return &item, nil
	}

	return nil, fmt.Errorf("no metadata found for run id %s", runId)
}

func getTagsQueryClause(tags map[string]string) ([]string, []interface{}, error) {
	clause := []string{}
	params := []interface{}{}

	for k, v := range tags {
		//TODO(rajatjindal): maybe check if this is a key in our db already?
		if strings.ContainsAny(k, " \t\n\r\"'`") {
			return nil, nil, fmt.Errorf("Invalid tag key")
		}

		clause = append(clause, fmt.Sprintf("json_extract(metadata.tags, '$.%s') = ?", k))
		params = append(params, v)
	}

	return clause, params, nil
}

func FetchAllRuns(ctx context.Context, filter types.CommonFilter) ([]*types.Summary, error) {
	conn := db()
	defer conn.Close()

	query := "select summary.* from summary, metadata where metadata.run_id = summary.run_id"
	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	rows, err := conn.Queryx(query, params...)
	if err != nil {
		return nil, err
	}

	items := []*types.Summary{}
	for rows.Next() {
		var item types.Summary
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

type Tag struct {
	Key            string   `json:"key" db:"tagkey"`
	Values         []string `json:"values"`
	CommaSepValues string   `json:"-" db:"comma_sep_values"` // used for scanning from db
}

func GetTagsForQuery(ctx context.Context, filter types.CommonFilter) ([]Tag, error) {
	conn := db()
	defer conn.Close()

	query := `WITH json_parsed AS (
    SELECT
        json_each.key AS key,
        json_each.value AS value
    FROM
        metadata,
        json_each(tags)
		WHERE
			repo = ?
)
SELECT
    key AS tagkey,
    GROUP_CONCAT(DISTINCT value) AS comma_sep_values
FROM
    json_parsed
GROUP BY
    key;
`
	rows, err := conn.QueryxContext(ctx, query, filter.Repo)
	if err != nil {
		return nil, err
	}

	items := []Tag{}
	for rows.Next() {
		var item Tag
		err = rows.StructScan(&item)
		if err != nil {
			return nil, errors.New(err.Error() + " scanning tag")
		}

		item.Values = strings.Split(item.CommaSepValues, ",")
		items = append(items, item)
	}

	return items, nil
}

// TODO(rajatjindal): do in a transaction
func IngestTestRun(ctx context.Context, metadata *types.Metadata, summary *types.Summary, suites []types.Suite) error {
	conn := db()
	defer conn.Close()

	tags, err := json.Marshal(metadata.Tags)
	if err != nil {
		return err
	}

	fmt.Println("ENTERING TAGS -> ", string(tags))
	_, err = conn.QueryxContext(ctx, "INSERT INTO metadata (run_id, repo, branch, commit_sha, job_name, format, link, tags, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", metadata.RunId, metadata.Repo, metadata.Branch, metadata.CommitSha, metadata.JobName, metadata.Format, metadata.Link, string(tags), metadata.CreatedAt)
	if err != nil {
		return err
	}

	_, err = conn.QueryxContext(ctx, "INSERT INTO summary (run_id, result, passed, failed, ignored, duration, created_at) values (?, ?, ?, ?, ?, ?, ?)", summary.RunId, summary.Result, summary.Passed, summary.Failed, summary.Ignored, summary.Duration, summary.CreatedAt)
	if err != nil {
		return err
	}

	// do batch update instead of one at a time
	suiteInsertStrings := []string{}
	suiteInsertParams := []interface{}{}
	testInsertStrings := []string{}
	testInsertParams := []interface{}{}
	total := 0
	for {
		for index, suite := range suites {
			total = total + 1
			suiteInsertStrings = append(suiteInsertStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			suiteInsertParams = append(suiteInsertParams, suite.RunId, index, suite.SuiteName, suite.Result, suite.Passed, suite.Failed, suite.Ignored, suite.Duration, suite.CreatedAt)

			for _, test := range suite.TestsTree {
				total = total + 1
				testInsertStrings = append(testInsertStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
				testInsertParams = append(testInsertParams, test.RunId, index, test.SuiteName, test.Name, test.Result, test.Duration, test.Logs, test.CreatedAt)
			}

			if total > 500 {
				fmt.Println("inserting ", total)
				// refresh the conn to hopefully fix "The stream has expired due to inactivity" error
				conn = db()
				suiteInsertSql := fmt.Sprintf("INSERT INTO suite_summary (run_id, suite_id, suite_name, result, passed, failed, ignored, duration, created_at) VALUES %s", strings.Join(suiteInsertStrings, ","))
				_, err = conn.QueryxContext(ctx, suiteInsertSql, suiteInsertParams...)
				if err != nil {
					return err
				}

				testsInsertSql := fmt.Sprintf("INSERT INTO tests (run_id, suite_id, suite_name, name, result, duration, logs, created_at) VALUES %s", strings.Join(testInsertStrings, ","))
				_, err = conn.QueryxContext(ctx, testsInsertSql, testInsertParams...)
				if err != nil {
					return err
				}

				total = 0
				testInsertStrings = []string{}
				testInsertParams = []interface{}{}
				suiteInsertStrings = []string{}
				suiteInsertParams = []interface{}{}
			}

		}

		conn = db()

		fmt.Println("inserting last ", total)
		//insert the remaining ones
		suiteInsertSql := fmt.Sprintf("INSERT INTO suite_summary (run_id, suite_id, suite_name, result, passed, failed, ignored, duration, created_at) VALUES %s", strings.Join(suiteInsertStrings, ","))
		_, err = conn.QueryxContext(ctx, suiteInsertSql, suiteInsertParams...)
		if err != nil {
			return err
		}

		testsInsertSql := fmt.Sprintf("INSERT INTO tests (run_id, suite_id, suite_name, name, result, duration, logs, created_at) VALUES %s", strings.Join(testInsertStrings, ","))
		_, err = conn.QueryxContext(ctx, testsInsertSql, testInsertParams...)
		if err != nil {
			return err
		}
		break
	}

	return nil
	// for index, suite := range suites {
	// 	suiteInsertStrings = append(suiteInsertStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	// 	suiteInsertParams = append(suiteInsertParams, suite.RunId, index, suite.SuiteName, suite.Result, suite.Passed, suite.Failed, suite.Ignored, suite.Duration, suite.CreatedAt)

	// 	for _, test := range suite.TestsTree {
	// 		testInsertStrings = append(testInsertStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
	// 		testInsertParams = append(testInsertParams, test.RunId, index, test.SuiteName, test.Name, test.Result, test.Duration, test.Logs, test.CreatedAt)
	// 	}
	// }

	// suiteInsertSql := fmt.Sprintf("INSERT INTO suite_summary (run_id, suite_id, suite_name, result, passed, failed, ignored, duration, created_at) VALUES %s", strings.Join(suiteInsertStrings, ","))
	// _, err = conn.QueryxContext(ctx, suiteInsertSql, suiteInsertParams...)
	// if err != nil {
	// 	return err
	// }

	// testsInsertSql := fmt.Sprintf("INSERT INTO tests (run_id, suite_id, suite_name, name, result, duration, logs, created_at) VALUES %s", strings.Join(testInsertStrings, ","))
	// _, err = conn.QueryxContext(ctx, testsInsertSql, testInsertParams...)
	// return err
}

func UpdateMetadata(ctx context.Context, runId, metadata string) error {
	conn := db()
	defer conn.Close()

	_, err := conn.QueryxContext(ctx, "UPDATE summary set metadata = ? WHERE run_id = ?;", metadata, runId)
	if err != nil {
		return err
	}

	return nil
}

func FetchTimeTrendsForSuite(ctx context.Context, filter *types.SuiteTrendsFilter) (*types.TimeTrendsData, error) {
	conn := db()
	defer conn.Close()

	clause, params, err := getClauseAndParamsFromFilter(filter)
	if err != nil {
		return nil, err
	}

	query := "select suite_summary.run_id as run_id, metadata.commit_sha as commit_sha, metadata.job_name, suite_summary.suite_name as suite_name, suite_summary.duration as duration from suite_summary, metadata WHERE suite_summary.run_id = metadata.run_id"
	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	perPage := filter.PerPage

	//defaults to defaultPerPage if filter.PerPage is not provided
	if perPage == 0 {
		perPage = defaultPerPage
	}

	query += " order by suite_summary.suite_name asc"
	if filter.Page > 0 {
		query += " offset ?"
		params = append(params, filter.Page*perPage)
	}

	if perPage > 0 {
		query += " limit ?"
		params = append(params, perPage)
	}

	fmt.Println(query)
	fmt.Println(params)
	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	items := []*types.SuiteTimeTrendEntry{}
	for rows.Next() {
		var item types.SuiteTimeTrendEntry
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	// TODO(rajatjindal): in the latest execution find the slowest 10 executions
	// then for all other runs prepare the execution times for those 10 suites only
	datamap := map[string]map[string]map[string]*types.SuiteTimeTrendEntry{}
	for _, item := range items {
		_, ok := datamap[item.CommitSha]
		if !ok {
			datamap[item.CommitSha] = map[string]map[string]*types.SuiteTimeTrendEntry{}
		}

		_, ok = datamap[item.CommitSha][item.RunId]
		if !ok {
			datamap[item.CommitSha][item.RunId] = map[string]*types.SuiteTimeTrendEntry{}
		}

		datamap[item.CommitSha][item.RunId][item.SuiteName] = item
	}

	sumdata := &types.TimeTrendsData{}
	suiteMap := map[string]types.Dataset{}
	for commitSha, abcMap := range datamap {
		sumdata.Labels = append(sumdata.Labels, []string{commitSha})
		for _, suitesMap := range abcMap {
			for suiteName, entry := range suitesMap {
				existing, ok := suiteMap[entry.RunId+suiteName]
				if !ok {
					existing = types.Dataset{
						Label:           entry.JobName + " : " + entry.SuiteName,
						Data:            []float64{},
						Stack:           entry.RunId,
						BackgroundColor: fmt.Sprintf("rgb(%d,%d,%d)", rand.IntN(255), rand.IntN(255), rand.IntN(255)),
					}
				}
				existing.Data = append(existing.Data, entry.Duration)

				suiteMap[entry.RunId+suiteName] = existing
			}
		}
	}

	for _, dataset := range suiteMap {
		sumdata.Datasets = append(sumdata.Datasets, dataset)
	}

	return sumdata, nil
}

func getClauseAndParamsFromFilter(filter *types.SuiteTrendsFilter) ([]string, []interface{}, error) {
	clause, params, err := getClauseAndParamsFromCommonFilter(filter.CommonFilter)
	if err != nil {
		return nil, nil, err
	}

	if filter.SuiteName != "" {
		parts := strings.Split(filter.SuiteName, " : ")
		if len(parts) > 1 {
			clause = append(clause, "metadata.job_name = ?")
			params = append(params, parts[0])

			clause = append(clause, "suite_summary.suite_name = ?")
			params = append(params, parts[1])
		} else {
			clause = append(clause, "suite_summary.suite_name = ?")
			params = append(params, parts[0])
		}
	}

	return clause, params, nil
}

func getClauseAndParamsFromCommonFilter(filter types.CommonFilter) ([]string, []interface{}, error) {
	if filter.Repo == "" {
		return nil, nil, fmt.Errorf("repo is mandatory")
	}

	clause, params, err := getTagsQueryClause(filter.Tags)
	if err != nil {
		return nil, nil, err
	}

	clause = append(clause, "metadata.repo = ?")
	params = append(params, filter.Repo)

	if filter.Branch != "" {
		clause = append(clause, "metadata.branch = ?")
		params = append(params, filter.Branch)
	}

	if filter.CommitSha != "" {
		clause = append(clause, "metadata.commit_sha = ?")
		params = append(params, filter.CommitSha)
	}

	return clause, params, nil
}

func FetchReliabilityTrends(ctx context.Context, filter types.CommonFilter) (*types.TimeTrendsData, error) {
	conn := db()
	defer conn.Close()

	// fetch data specific for the commit
	if filter.CommitSha != "" {
		return fetchReliabilityTrendsForCommit(ctx, filter)
	}

	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	fmt.Printf("clause: %v, params: %v\n", clause, params)
	query := `
SELECT 
	metadata.commit_sha as commit_sha,
	metadata.created_at as created_at,
	sum(summary.passed) as passed, 
	sum(summary.failed) as failed, 
	sum(summary.ignored) as ignored 
FROM
 summary, 
 metadata 
WHERE
  metadata.run_id = summary.run_id`

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	query += ` GROUP BY metadata.commit_sha ORDER BY created_at`

	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	type TempData struct {
		CommitSha string `db:"commit_sha"`
		CreatedAt string `db:"created_at"`
		Passed    int    `db:"passed"`
		Failed    int    `db:"failed"`
		Ignored   int    `db:"ignored"`
	}

	sumdata := &types.TimeTrendsData{
		Datasets: []types.Dataset{
			{
				Label:           "passed",
				Data:            []float64{},
				BackgroundColor: "#34e8bd",
				Hidden:          true,
			},
			{
				Label:           "failed",
				Data:            []float64{},
				BackgroundColor: "#CE5050",
			},
			{
				Label:           "skipped",
				Data:            []float64{},
				BackgroundColor: "#9ca3af",
			},
		},
	}

	items := []*TempData{}
	for rows.Next() {
		var item TempData
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		//(TODO): TEMP WORKAROUND FOR TESTS. IN ACTUAL DATE WLL BE PUSHED
		if item.CreatedAt == "" {
			item.CreatedAt = time.Now().Format(time.RFC3339)
		}
		createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			return nil, err
		}

		sumdata.Labels = append(sumdata.Labels, []string{item.CommitSha, createdAt.Format("01/02 15:04")})
		sumdata.Ids = append(sumdata.Ids, item.CommitSha) // helps with creating links
		sumdata.Datasets[0].Data = append(sumdata.Datasets[0].Data, float64(item.Passed))
		sumdata.Datasets[1].Data = append(sumdata.Datasets[1].Data, float64(item.Failed))
		sumdata.Datasets[2].Data = append(sumdata.Datasets[2].Data, float64(item.Ignored))
		items = append(items, &item)
	}

	/*
		labels = [a, b]
		datasets = [
			{
				label: passed,
				data: [a, b]
				stack: "",
				color: ""
			}
		]


	*/

	return sumdata, nil
}

func fetchReliabilityTrendsForCommit(ctx context.Context, filter types.CommonFilter) (*types.TimeTrendsData, error) {
	conn := db()
	defer conn.Close()

	if filter.CommitSha == "" {
		return nil, fmt.Errorf("commitSha is required")
	}

	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	fmt.Printf("clause: %v, params: %v\n", clause, params)
	query := `
SELECT 
	metadata.run_id as run_id,
	metadata.tags as tags,
	sum(summary.passed) as passed, 
	sum(summary.failed) as failed, 
	sum(summary.ignored) as ignored 
FROM
 summary, 
 metadata 
WHERE
  metadata.run_id = summary.run_id`

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	query += ` GROUP BY metadata.run_id`

	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	type TempData struct {
		RunId   string `db:"run_id"`
		Tags    string `db:"tags"`
		Passed  int    `db:"passed"`
		Failed  int    `db:"failed"`
		Ignored int    `db:"ignored"`
	}

	sumdata := &types.TimeTrendsData{
		Datasets: []types.Dataset{
			{
				Label:           "passed",
				Data:            []float64{},
				BackgroundColor: "#34e8bd",
				Hidden:          true,
			},
			{
				Label:           "failed",
				Data:            []float64{},
				BackgroundColor: "#CE5050",
			},
			{
				Label:           "skipped",
				Data:            []float64{},
				BackgroundColor: "#9ca3af",
			},
		},
	}

	items := []*TempData{}
	for rows.Next() {
		var item TempData
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		tags := map[string]string{}
		err = json.Unmarshal([]byte(item.Tags), &tags)
		if err != nil {
			return nil, err
		}

		sumdata.Labels = append(sumdata.Labels, []string{tags["engine_version"], tags["workflow"], tags["job"]})
		sumdata.Ids = append(sumdata.Ids, item.RunId) // helps with creating links
		sumdata.Datasets[0].Data = append(sumdata.Datasets[0].Data, float64(item.Passed))
		sumdata.Datasets[1].Data = append(sumdata.Datasets[1].Data, float64(item.Failed))
		sumdata.Datasets[2].Data = append(sumdata.Datasets[2].Data, float64(item.Ignored))
		items = append(items, &item)
	}

	return sumdata, nil
}

func FetchTopNSlowestTestSuites(ctx context.Context, filter types.CommonFilter) ([]TopNSlowestTestSuites, error) {
	conn := db()
	defer conn.Close()

	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	query := `WITH ranked_suites AS (
		SELECT
			suite_summary.suite_name,
			suite_summary.duration,
			ROW_NUMBER() OVER (PARTITION BY suite_summary.suite_name ORDER BY suite_summary.duration) AS rank,
			COUNT(*) OVER (PARTITION BY suite_summary.suite_name) AS total_count
		FROM 
			suite_summary, metadata
		WHERE
			suite_summary.run_id = metadata.run_id`

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	query += `),
	p90_positions AS (
		SELECT 
			suite_name,
			CEIL(0.9 * total_count) AS p90_pos
		FROM ranked_suites
		GROUP BY suite_name, total_count
	),
	p90_scores AS (
		SELECT 
			rs.suite_name,
			rs.duration AS p90_duration
		FROM ranked_suites rs
		JOIN p90_positions pp
			ON rs.suite_name = pp.suite_name AND rs.rank = pp.p90_pos
	)
	SELECT 
		suite_name as suite_name,
		p90_duration as duration
	FROM p90_scores
	ORDER BY p90_duration DESC
	LIMIT 10;`

	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	items := []TopNSlowestTestSuites{}
	for rows.Next() {
		var item TopNSlowestTestSuites
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func FetchTopNFlakyTestSuites(ctx context.Context, filter types.CommonFilter) ([]TopNFlakyTestSuites, error) {
	conn := db()
	defer conn.Close()

	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
    suite_name,
    SUM(failed) * 100.0 / (SUM(failed) + SUM(passed) + SUM(ignored)) AS failure_percentage
FROM 
    suite_summary, metadata
WHERE
			suite_summary.run_id = metadata.run_id`

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	query += ` 
GROUP BY 
    suite_name
HAVING 
    SUM(failed) + SUM(passed) + SUM(ignored) > 0  -- Exclude suites with no executed tests
ORDER BY 
    failure_percentage DESC
LIMIT 10;`

	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	items := []TopNFlakyTestSuites{}
	for rows.Next() {
		var item TopNFlakyTestSuites
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

type TopNSlowestTestSuites struct {
	SuiteName string `db:"suite_name"`
	Duration  string `db:"duration"`
}

type TopNFlakyTestSuites struct {
	SuiteName         string `db:"suite_name"`
	FailurePercentage string `db:"failure_percentage"`
}
