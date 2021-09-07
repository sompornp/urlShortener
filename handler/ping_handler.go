package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type PingHandler struct {
}

func (h *PingHandler) ConfigRoute(r *mux.Router) {
	r.HandleFunc("/ping", h.Handler)
}

func (h *PingHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}
