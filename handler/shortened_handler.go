package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sompornp/urlShortener/constant"
	"github.com/sompornp/urlShortener/db"
	"github.com/sompornp/urlShortener/model"
	"github.com/sompornp/urlShortener/utils"
	"net/http"
	"net/url"
	"time"
)

type ShortenedHandler struct {
	DB         *db.DB
	RefererUrl *url.URL
	Blacklists []model.Blacklist
	Token      string
}

func (h *ShortenedHandler) ConfigRoute(r *mux.Router, token string) {
	r.HandleFunc("/encode", h.Encode).Methods("POST")
	r.HandleFunc("/{shortLink}", h.ShortLink)
}

func (h *ShortenedHandler) ShortLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["shortLink"]

	links := h.DB.FindShortLink(&db.DBCriteria{Id: shortLink})

	if len(links) == 0 {
		code := http.StatusNotFound
		http.Error(w, http.StatusText(code), code)
		return
	}

	link := links[0]

	if time.Now().Unix() > link.Expire {
		code := http.StatusGone
		http.Error(w, http.StatusText(code), code)
		return
	}

	h.DB.UpdateShortLink(link.Id)

	http.Redirect(w, r, link.TargetUrl, http.StatusFound)
}

func (h *ShortenedHandler) Encode(w http.ResponseWriter, r *http.Request) {
	type encodeInput struct {
		Url    string
		Expire int64
	}

	type EncodeResponse struct {
		Success bool
		Url     string
	}

	var input encodeInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}
	defer r.Body.Close()

	targetUrl, urlErr := url.Parse(input.Url)
	if urlErr != nil || targetUrl.Scheme == "" || targetUrl.Host == "" {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}

	if input.Expire < 0 || input.Expire > constant.MaxUnixTimestamp {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}

	if utils.ContainBlacklist(h.Blacklists, targetUrl.String()) {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code)+" - contains blacklist", code)
		return
	}

	id, _ := gonanoid.New(6)

	for len(h.DB.FindShortLink(&db.DBCriteria{Id: id})) > 0 {
		id, _ = gonanoid.New(6)
	}

	shortUrl := url.URL{
		Scheme: h.RefererUrl.Scheme,
		Host:   h.RefererUrl.Host,
		Path:   id}

	h.DB.InsertShortLink(model.ShortLink{
		Id:        id,
		ShortUrl:  shortUrl.String(),
		TargetUrl: targetUrl.String(),
		Cnt:       0,
		Expire:    input.Expire,
	})

	response, _ := json.Marshal(EncodeResponse{
		Success: true,
		Url:     shortUrl.String(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}