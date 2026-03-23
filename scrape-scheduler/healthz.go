package scrape_scheduler

import (
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
