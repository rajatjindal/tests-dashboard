package main

import (
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/rust-test-reporter/backend/pkg/api"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		router := api.New()
		router.ServeHTTP(w, r)
	})
}

func main() {}
