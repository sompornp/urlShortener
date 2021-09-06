package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type PinHandler struct {
}

func (h *PinHandler) ConfigRoute(r *mux.Router) {
	r.HandleFunc("/ping", h.Handler)
}

func (h *PinHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}
