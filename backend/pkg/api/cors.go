package api

import (
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

func cors(router *spinhttp.Router) {
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", "*")
		header.Set("Access-Control-Allow-Origin", "*")

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
}
