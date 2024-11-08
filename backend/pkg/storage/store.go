package storage

import (
	"context"
	"fmt"

	"github.com/fermyon/spin-go-sdk/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

func db() *sqlx.DB {
	conn := sqlite.Open("default")
	return sqlx.NewDb(conn, "sqlite")
}

func FetchAllTests(ctx context.Context) ([]*types.Test, error) {
	conn := db()
	if conn == nil {
		return nil, fmt.Errorf("failed to read db")
	}
	defer conn.Close()

	rows, err := conn.Queryx("select * from tests;")
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

func FetchHistoryForLogLine(ctx context.Context, logLine string) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from tests where lower(logs) LIKE ?", "%"+logLine+"%")
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

func FetchHistoryForTestcase(ctx context.Context, testname string) ([]*types.Test, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select * from tests where lower(name) LIKE ?", "%"+testname+"%")
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
			run_id as run_id,
			suite_id as suite_id, 
			suite_name as suite_name, 
			SUM(CASE WHEN result = 'ok' THEN 1 ELSE 0 END) AS passed, 
			SUM(CASE WHEN result = 'failed' THEN 1 ELSE 0 END) AS failed,
			SUM(CASE WHEN result = 'ignored' THEN 1 ELSE 0 END) AS ignored 
		from 
			tests 
		where 
			run_id = ? group by suite_id;`, runId)
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

func FetchAllRuns(ctx context.Context, repo string) ([]*types.Summary, error) {
	conn := db()
	defer conn.Close()

	rows, err := conn.Queryx("select summary.* from summary, metadata where metadata.repo = ? AND metadata.run_id = summary.run_id;", repo)
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

func IngestTestRun(ctx context.Context, metadata *types.Metadata, summary *types.Summary, suites []types.Suite) error {
	conn := db()
	defer conn.Close()

	_, err := conn.QueryxContext(ctx, "INSERT INTO metadata values (?, ?, ?, ?, ?, ?, ?)", metadata.RunId, metadata.Repo, metadata.Branch, metadata.Format, metadata.Link, metadata.Tags, metadata.CreatedAt)
	if err != nil {
		return err
	}

	_, err = conn.QueryxContext(ctx, "INSERT INTO summary values (?, ?, ?, ?, ?, ?, ?)", summary.RunId, summary.Result, summary.Passed, summary.Failed, summary.Ignored, summary.Duration, summary.CreatedAt)
	if err != nil {
		return err
	}

	for index, suite := range suites {
		//TODO(rajatjindal): WHY IS INDEX STARTING FROM 1 HERE?
		nindex := index - 1
		for _, test := range suite.TestsTree {
			fmt.Println("index is -> %d", index)
			_, err := conn.QueryxContext(ctx, "INSERT INTO tests values (?, ?, ?, ?, ?, ?, ?, ?)", test.RunId, fmt.Sprintf("%d", nindex), test.SuiteName, test.Name, test.Result, test.Duration, test.Logs, test.CreatedAt)
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
