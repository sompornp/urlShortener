package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sompornp/urlShortener/db"
	"github.com/sompornp/urlShortener/model"
	"net/http"
	"net/url"
)

type AdminHandler struct {
	DB         *db.DB
	RefererUrl *url.URL
	Blacklists []model.Blacklist
	Token      string
}

func (h *AdminHandler) ConfigRoute(r *mux.Router, token string) {
	r.HandleFunc("/admin", middleware(http.HandlerFunc(h.List), token)).Methods("GET")
	r.HandleFunc("/admin", middleware(http.HandlerFunc(h.Delete), token)).Methods("DELETE")
}

func (h *AdminHandler) List(w http.ResponseWriter, r *http.Request) {
	var cri = db.DBCriteria{
		Id:      r.URL.Query().Get("Shortcode"),
		Keyword: r.URL.Query().Get("Keyword"),
	}

	shortLinks := h.DB.FindShortLink(&cri)

	response, _ := json.Marshal(shortLinks)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	type deleteInput struct {
		Shortcode string
	}

	var input deleteInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}
	defer r.Body.Close()

	if input.Shortcode == "" {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}

	if len(h.DB.FindShortLink(&db.DBCriteria{Id: input.Shortcode})) > 0 {
		h.DB.DeleteShortLink(input.Shortcode)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		code := http.StatusGone
		http.Error(w, http.StatusText(code), code)
	}
}

func middleware(next http.Handler, token string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inputToken := r.Header.Get("Authorization")
		if token != inputToken {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		next.ServeHTTP(w, r)
	})
}
