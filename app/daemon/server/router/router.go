package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", healthcheck)
	return r
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok"})
}
