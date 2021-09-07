package handler

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingConfigure(t *testing.T) {
	handler := PingHandler{}
	router := mux.NewRouter()
	handler.ConfigRoute(router)
}

func TestPingHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	handler := PingHandler{}

	handler.Handler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []byte("pong"), w.Body.Bytes())
}
