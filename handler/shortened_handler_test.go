package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/sompornp/urlShortener/constant"
	"github.com/sompornp/urlShortener/db"
	"github.com/sompornp/urlShortener/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestShortLinkConfigure(t *testing.T) {
	handler := ShortenedHandler{}
	router := mux.NewRouter()
	handler.ConfigRoute(router, "abc")
}

func TestShortLinkEncodeHandler(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectQuery("SELECT (.+) FROM shortLink where 1 = 1 and id = (.+)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "shortUrl", "targetUrl", "cnt", "expire"}))

	mock.ExpectPrepare("INSERT INTO shortlink").
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r, _ := http.NewRequest("POST", "/encode",
		strings.NewReader(`{"Url": "http://www.ku.ac.th","Expire": 1630990800}`))
	w := httptest.NewRecorder()

	handler := ShortenedHandler{
		DB:         &db.DB{DB: mockDb},
		RefererUrl: &url.URL{Scheme: "http", Host: "localhost:8080"},
		Blacklists: []model.Blacklist{},
		Token:      "abc",
	}

	handler.Encode(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body.Bytes())
}

func TestShortLinkHandler(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectQuery("SELECT (.+) FROM shortLink where 1 = 1 and id = (.+)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "shortUrl", "targetUrl", "cnt", "expire"}).
			AddRow("1", "http://localhost:8080/abcdef", "http://www.ku.ac.th", "0", constant.MaxUnixTimestamp))

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE shortLink SET cnt = (.+) where id = (.+)").
		WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r, _ := http.NewRequest("GET", "/abcdef", nil)
	r = mux.SetURLVars(r, map[string]string{
		"shortLink": "abcdef",
	})

	w := httptest.NewRecorder()

	handler := ShortenedHandler{
		DB:         &db.DB{DB: mockDb},
		RefererUrl: &url.URL{Scheme: "http", Host: "localhost:8080"},
		Blacklists: []model.Blacklist{},
		Token:      "abc",
	}

	handler.ShortLink(w, r)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, strings.TrimSpace(`<a href="http://www.ku.ac.th">Found</a>.`), strings.TrimSpace(string(w.Body.Bytes())))
}
