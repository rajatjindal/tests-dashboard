package api

import (
	"encoding/json"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

func New() *spinhttp.Router {
	router := spinhttp.NewRouter()
	cors(router)

	router.GET("/api/healthz", func(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
		w.WriteHeader(http.StatusOK)
	})

	router.POST("/api/schema", createSchema)

	router.GET("/api/tags", getTags)
	router.POST("/api/run/:runId", ingestTestRun)

	router.GET("/api/trends/reliability", fetchReliabilityTrends)

	router.GET("/api/runs", fetchAllRuns)
	router.GET("/api/runs/:runId/metadata", fetchMetadataForRun)
	router.GET("/api/runs/:runId/summary", fetchSummaryForRun)
	router.GET("/api/runs/:runId/suites", fetchSuitesForRun)
	router.GET("/api/runs/:runId/suites-summary", fetchSuiteSummaryForRunId)
	router.GET("/api/runs/:runId/suites/:suiteId/tests", fetchTestsForRunIdAndSuite)

	// reports
	router.GET("/api/reports/top-n-slowest-tests", fetchTopNSlowestTestSuites)
	router.GET("/api/reports/top-n-flaky-tests", fetchTopNFlakyTestSuites)

	// time trends
	router.GET("/api/trends/suites/time", fetchTimeTrendsForSuite)

	router.GET("/api/history/log", fetchHistoryForLogLine)
	router.GET("/api/history/test", fetchHistoryForTestcase)

	return router
}

func prettyPrint(w http.ResponseWriter, data interface{}) {
	header := w.Header()
	header.Set("Access-Control-Allow-Methods", "*")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}
