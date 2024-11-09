package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/gojson"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/junit"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/rustjson"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/storage"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

func fetchAllRuns(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	summary, err := storage.FetchAllRuns(r.Context(), repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func fetchMetadataForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	summary, err := storage.FetchMetadataForRun(r.Context(), params.ByName("runId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func fetchSummaryForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	summary, err := storage.FetchSummaryForRun(r.Context(), params.ByName("runId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func updateMetadata(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = storage.UpdateMetadata(r.Context(), runId, string(raw))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ingestTestRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	fmt.Printf("run id %s\n", runId)

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var metadata types.Metadata
	var rawResults []byte

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("part %#v\n", part)

		_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if params["name"] == "metadata" {
			raw, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("metadata is %s\n", string(raw))
			err = json.Unmarshal(raw, &metadata)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if params["name"] == "results" {
			rawResults, err = io.ReadAll(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var summary *types.Summary
	var suites []types.Suite
	if metadata.Format == "junit" {
		summary, suites, err = junit.Ingest(runId, rawResults)
	} else if metadata.Format == "gojson" {
		summary, suites, err = gojson.Ingest(runId, rawResults)
	} else {
		summary, suites, err = rustjson.Ingest(runId, rawResults)
	}

	err = storage.IngestTestRun(r.Context(), &metadata, summary, suites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func fetchSuiteSummaryForRunId(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")

	tests, err := storage.FetchSuiteSummaryForRunId(r.Context(), runId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchTestsForRunIdAndSuite(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	suiteId := params.ByName("suiteId")

	tests, err := storage.FetchTestsByRunIdAndSuite(r.Context(), runId, suiteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchHistoryForLogLine(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	logLine := r.URL.Query().Get("logLine")

	tests, err := storage.FetchHistoryForLogLine(r.Context(), logLine)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchHistoryForTestcase(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	logLine := r.URL.Query().Get("name")

	tests, err := storage.FetchHistoryForTestcase(r.Context(), logLine)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchSuitesForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	suites, err := storage.FetchSuitesForRun(r.Context(), runId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, suites)
}

func fetchAllTests(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	tests, err := storage.FetchAllTests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}
