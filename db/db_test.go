package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sompornp/urlShortener/constant"
	"github.com/sompornp/urlShortener/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertShortLink(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectPrepare("INSERT INTO shortlink").
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{
		DB: mockDb,
	}

	db.InsertShortLink(model.ShortLink{
		Id:        "abc",
		ShortUrl:  "http://localhost:8080/abcdef",
		TargetUrl: "http://www.ku.ac.th",
		Cnt:       0,
		Expire:    constant.MaxUnixTimestamp,
	})
}

func TestFindShortLink(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectQuery("SELECT (.+) FROM shortLink").
		WillReturnRows(sqlmock.NewRows([]string{"id", "shortUrl", "targetUrl", "cnt", "expire"}).
			AddRow("1", "http://localhost:8080/abcdef", "http://www.ku.ac.th", "0", constant.MaxUnixTimestamp))

	db := DB{
		DB: mockDb,
	}

	shortlinks := db.FindShortLink(&DBCriteria{})

	assert.Equal(t, 1, len(shortlinks))
	assert.Equal(t, "1", shortlinks[0].Id)
	assert.Equal(t, "http://localhost:8080/abcdef", shortlinks[0].ShortUrl)
}

func TestDeleteShortLink(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectPrepare("DELETE FROM shortlink where id = (.+)").
		ExpectExec().WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{
		DB: mockDb,
	}

	db.DeleteShortLink("abcdef")
}

func TestUpdateShortLink(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE shortLink SET cnt = (.+) where id = (.+)").
		WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	db := DB{
		DB: mockDb,
	}

	db.UpdateShortLink("abcdef")
}
