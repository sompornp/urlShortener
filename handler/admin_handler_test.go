package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/sompornp/urlShortener/db"
	"github.com/sompornp/urlShortener/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAdminConfigure(t *testing.T) {
	handler := AdminHandler{}
	router := mux.NewRouter()
	handler.ConfigRoute(router, "abc")
}

func TestAdminListHandler(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectQuery("SELECT (.+) FROM shortLink").
		WillReturnRows(sqlmock.NewRows([]string{"id", "shortUrl", "targetUrl", "cnt", "expire"}).
			AddRow("1", "http://localhost:8080/abcdef", "http://www.ku.ac.th", "0", 1631011785))

	r, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	handler := AdminHandler{
		DB:         &db.DB{DB: mockDb},
		RefererUrl: &url.URL{Scheme: "http", Host: "localhost:8080"},
		Blacklists: []model.Blacklist{},
		Token:      "abc",
	}

	handler.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `
[{"Id":"1","ShortUrl":"http://localhost:8080/abcdef","TargetUrl":"http://www.ku.ac.th","Hits":0,"Expire":1631011785}]
`, string(w.Body.Bytes()))
}

func TestAdminDeleteHandler(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectQuery("SELECT (.+) FROM shortLink where (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id", "shortUrl", "targetUrl", "cnt", "expire"}).
			AddRow("1", "http://localhost:8080/abcdef", "http://www.ku.ac.th", "0", 1631011785))

	mock.ExpectPrepare("DELETE FROM shortlink where id = (.+)").
		ExpectExec().WithArgs("qHMLYY").WillReturnResult(sqlmock.NewResult(1, 1))

	r, _ := http.NewRequest("DELETE", "/admin", strings.NewReader(`{"Shortcode": "qHMLYY"}`))
	w := httptest.NewRecorder()

	handler := AdminHandler{
		DB:         &db.DB{DB: mockDb},
		RefererUrl: &url.URL{Scheme: "http", Host: "localhost:8080"},
		Blacklists: []model.Blacklist{},
		Token:      "abc",
	}

	handler.Delete(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", string(w.Body.Bytes()))
}
