package api

import (
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/storage"
)

func createSchema(w http.ResponseWriter, r *http.Request, params spinhttp.Params) {
	err := storage.CreateSchema(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
