package main

import (
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/api"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		router := api.New()
		router.ServeHTTP(w, r)
	})
}

func main() {}
