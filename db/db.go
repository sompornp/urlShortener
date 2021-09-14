package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sompornp/urlShortener/model"
	"log"
	"os"
)

type DBCriteria struct {
	Id      string
	Keyword string
}

type DB struct {
	DB *sql.DB
}

func New(dbFile string) *DB {

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

	db, _ := sql.Open("sqlite3", dbFileName)
	createTableIfNotExist(db)

	return &DB{DB: db}
}

func createTableIfNotExist(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS shortLink (
		"id" TEXT NOT NULL PRIMARY KEY,		
		"shortUrl" TEXT NOT NULL,
		"targetUrl" TEXT NOT NULL,
		"cnt" INTEGER NOT NULL,
		"expire" INTEGER NOT NULL	
	  );`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func (d *DB) InsertShortLink(link model.ShortLink) error {
	insertSQL := `INSERT INTO shortlink(id, shortUrl, targetUrl, cnt, expire) VALUES (?, ?, ?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(link.Id, link.ShortUrl, link.TargetUrl, link.Cnt, link.Expire)
	return err
}

func (d *DB) UpdateShortLink(id string) {
	ctx := context.Background()
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE shortLink SET cnt = cnt+1 where id = ?`, id)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Fatalln(err.Error())
	}
}

func (d *DB) DeleteShortLink(id string) {
	deleteSQL := "DELETE FROM shortlink where id = ?"
	statement, err := d.DB.Prepare(deleteSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (d *DB) FindShortLink(cri *DBCriteria) []model.ShortLink {
	var args []interface{}
	where := " where 1 = 1"

	if len(cri.Id) > 0 {
		where += " and id = ?"
		args = append(args, cri.Id)
	}

	if len(cri.Keyword) > 0 {
		where += " and targetUrl like ?"
		args = append(args, fmt.Sprintf("%%%s%%", cri.Keyword))
	}

	row, err := d.DB.Query("SELECT id, shortUrl, targetUrl, cnt, expire FROM shortLink "+where, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	shortLinks := []model.ShortLink{}
	for row.Next() {
		var sl model.ShortLink
		row.Scan(&sl.Id, &sl.ShortUrl, &sl.TargetUrl, &sl.Cnt, &sl.Expire)
		shortLinks = append(shortLinks, sl)
	}
	return shortLinks
}
