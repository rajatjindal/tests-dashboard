package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/gojson"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/junit"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/parser/rustjson"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/storage"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
)

func getCommonFilterFromRequest(r *http.Request) types.CommonFilter {
	fmt.Println("query -> ", r.URL.Query())
	return types.CommonFilter{
		Repo:      r.URL.Query().Get("repo"),
		Branch:    r.URL.Query().Get("branch"),
		CommitSha: r.URL.Query().Get("commitSha"),
		Tags:      getTagsFromQuery(r),
	}
}

func getTagsFromQuery(r *http.Request) map[string]string {
	tags := map[string]string{}
	for k, v := range r.URL.Query() {
		if after, found := strings.CutPrefix(k, "tag-"); found && len(v) > 0 {
			tags[after] = v[0]
		}
	}

	return tags
}

func fetchReliabilityTrends(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	filter := getCommonFilterFromRequest(r)
	fmt.Printf("reliability: %#v", filter)
	summary, err := storage.FetchReliabilityTrends(r.Context(), filter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func fetchAllRuns(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	filter := getCommonFilterFromRequest(r)
	summary, err := storage.FetchAllRuns(r.Context(), filter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func fetchMetadataForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	summary, err := storage.FetchMetadataForRun(r.Context(), params.ByName("runId"))
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func fetchSummaryForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	summary, err := storage.FetchSummaryForRun(r.Context(), params.ByName("runId"))
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, summary)
}

func updateMetadata(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = storage.UpdateMetadata(r.Context(), runId, string(raw))
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTags(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	filter := getCommonFilterFromRequest(r)
	tags, err := storage.GetTagsForQuery(r.Context(), filter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tags)
}

func ingestTestRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	fmt.Printf("run id %s\n", runId)

	reader, err := r.MultipartReader()
	if err != nil {
		fmt.Println("ERROR ", err)
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
			fmt.Println("ERROR ", err)
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
				fmt.Println("ERROR ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if params["name"] == "results" {
			rawResults, err = io.ReadAll(part)
			if err != nil {
				fmt.Println("ERROR ", err)
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
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = storage.IngestTestRun(r.Context(), &metadata, summary, suites)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func fetchSuiteSummaryForRunId(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")

	tests, err := storage.FetchSuiteSummaryForRunId(r.Context(), runId)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchTopNSlowestTestSuites(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	filter := getCommonFilterFromRequest(r)
	items, err := storage.FetchTopNSlowestTestSuites(r.Context(), filter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, items)
}

func fetchTopNFlakyTestSuites(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	filter := getCommonFilterFromRequest(r)
	items, err := storage.FetchTopNFlakyTestSuites(r.Context(), filter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, items)
}

func fetchTestsForRunIdAndSuite(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	suiteId := params.ByName("suiteId")

	tests, err := storage.FetchTestsByRunIdAndSuite(r.Context(), runId, suiteId)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchTimeTrendsForSuite(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	commonFilter := getCommonFilterFromRequest(r)
	commonFilter.Page = 0
	commonFilter.PerPage = -1

	suiteName := r.URL.Query().Get("suiteName")
	trends, err := storage.FetchTimeTrendsForSuite(r.Context(), &types.SuiteTrendsFilter{
		SuiteName:    suiteName,
		CommonFilter: commonFilter,
	})
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, trends)
}

func fetchHistoryForLogLine(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	commonFilter := getCommonFilterFromRequest(r)

	logLine := r.URL.Query().Get("logLine")

	tests, err := storage.FetchHistoryForLogLine(r.Context(), logLine, commonFilter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchHistoryForTestcase(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	commonFilter := getCommonFilterFromRequest(r)

	testcase := r.URL.Query().Get("name")

	tests, err := storage.FetchHistoryForTestcase(r.Context(), testcase, commonFilter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}

func fetchSuitesForRun(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	runId := params.ByName("runId")
	suites, err := storage.FetchSuitesForRun(r.Context(), runId)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, suites)
}

func fetchAllTests(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "repo is required", http.StatusBadRequest)
		return
	}

	commonFilter := getCommonFilterFromRequest(r)

	tests, err := storage.FetchAllTests(r.Context(), commonFilter)
	if err != nil {
		fmt.Println("ERROR ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prettyPrint(w, tests)
}
