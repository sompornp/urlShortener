package main

import (
	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	db2 "github.com/sompornp/urlShortener/db"
	"github.com/sompornp/urlShortener/handler"
	"github.com/sompornp/urlShortener/utils"
	"os"
)
import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbFileName := os.Getenv("DB_FILE")
	if dbFileName == "" {
		log.Fatal("DB_FILE is required")
	}

	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		_, err := os.Create(dbFileName)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	db := db2.New(dbFileName)
	defer db.DB.Close()

	referer := os.Getenv("REFERER")
	if referer == "" {
		referer = "http://localhost:8080"
	}

	refererUrl, err := url.Parse(referer)
	if err != nil || refererUrl.Scheme == "" || refererUrl.Host == "" {
		log.Fatal(err.Error())
	}

	token := os.Getenv("ADMIN_TOKEN")
	if token == "" {
		t, _ := gonanoid.New(12)
		token = t
	}
	log.Println("Token: ", token)

	blacklists := utils.BuildBlacklistUrl(os.Getenv("BLACKLIST_URLS"))

	r := mux.NewRouter()

	(&handler.PinHandler{}).ConfigRoute(r)

	(&handler.AdminHandler{
		DB:         db,
		RefererUrl: refererUrl,
		Blacklists: blacklists,
		Token:      token,
	}).ConfigRoute(r, token)

	(&handler.ShortenedHandler{
		DB:         db,
		RefererUrl: refererUrl,
		Blacklists: blacklists,
		Token:      token,
	}).ConfigRoute(r, token)

	log.Fatal(http.ListenAndServe(":8080", r))
}
