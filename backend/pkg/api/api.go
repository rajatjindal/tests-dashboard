package api

import (
	"encoding/json"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
)

func New() *spinhttp.Router {
	router := spinhttp.NewRouter()
	cors(router)

	router.POST("/api/run/:runId", ingestTestRun)

	router.GET("/api/runs", fetchAllRuns)
	router.GET("/api/runs/:runId/metadata", fetchMetadataForRun)
	router.GET("/api/runs/:runId/summary", fetchSummaryForRun)
	router.GET("/api/runs/:runId/suites", fetchSuitesForRun)
	router.GET("/api/runs/:runId/suites/:suiteId/tests", fetchTestsForRunIdAndSuite)

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
