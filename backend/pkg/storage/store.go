package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

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

	query := "select tests.* from tests, metadata where metadata.run_id = tests.run_id AND lower(tests.name) LIKE ?"
	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	if len(clause) > 0 {
		query += " AND "
		query += strings.Join(clause, " AND ")
	}

	rows, err := conn.Queryx(query, append([]interface{}{"%" + testname + "%"}, params...)...)
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

func FetchTestsByRunIdAndSuite(ctx context.Context, runId, suite string) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from tests where run_id = ? AND suite_id = ?;", runId, suite)
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

	for index, suite := range suites {
		_, err := conn.QueryxContext(ctx, "INSERT INTO suite_summary (run_id, suite_id, suite_name, result, passed, failed, ignored, duration, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", suite.RunId, index, suite.SuiteName, suite.Result, suite.Passed, suite.Failed, suite.Ignored, suite.Duration, suite.CreatedAt)
		if err != nil {
			return err
		}

		for _, test := range suite.TestsTree {
			// fmt.Println("index is -> %d", index)
			_, err := conn.QueryxContext(ctx, "INSERT INTO tests (run_id, suite_id, suite_name, name, result, duration, logs, created_at) values (?, ?, ?, ?, ?, ?, ?, ?)", test.RunId, fmt.Sprintf("%d", index), test.SuiteName, test.Name, test.Result, test.Duration, test.Logs, test.CreatedAt)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
		sumdata.Labels = append(sumdata.Labels, commitSha)
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

	clause, params, err := getClauseAndParamsFromCommonFilter(filter)
	if err != nil {
		return nil, err
	}

	fmt.Printf("clause: %v, params: %v\n", clause, params)
	query := `
SELECT 
	metadata.commit_sha as commit_sha,
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

	query += ` GROUP BY metadata.commit_sha`

	rows, err := conn.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	type TempData struct {
		CommitSha string `db:"commit_sha"`
		Passed    int    `db:"passed"`
		Failed    int    `db:"failed"`
		Ignored   int    `db:"ignored"`
	}

	sumdata := &types.TimeTrendsData{
		Datasets: []types.Dataset{
			{
				Label:           "passed",
				Data:            []float64{},
				BackgroundColor: "green",
			},
			{
				Label:           "failed",
				Data:            []float64{},
				BackgroundColor: "red",
			},
			{
				Label:           "ignored",
				Data:            []float64{},
				BackgroundColor: "gray",
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

		sumdata.Labels = append(sumdata.Labels, item.CommitSha)
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
