package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goldencoderam/prg/database"
)

type Request struct {
	PersonId   int
	StatusType string
	Date       string
	Country    string
}

func (r Request) NewFakeRecord() (database.Insertable, error) {
	var personId int
	var statusType string

	database.ExecuteRowQuery("SELECT id FROM person ORDER BY random() LIMIT 1", &personId)
	database.ExecuteRowQuery("SELECT status_type FROM status_type ORDER BY random() LIMIT 1", &statusType)

	return Request{
		PersonId:   personId,
		StatusType: statusType,
		Date:       gofakeit.DateRange(
            time.Date(2020, 1,1,1,1,1,1, time.Now().Location()),
            time.Date(2022, 1,1,1,1,1,1, time.Now().Location()),
        ).Format("2006-01-02 15:04:05"),
		Country:    gofakeit.Country(),
	}, nil
}

func (r Request) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO request (person_id, status_type, date, country)
        VALUES (%d, '%s', '%s', '%s')`,
		r.PersonId,
		r.StatusType,
		r.Date,
		strings.ReplaceAll(r.Country, "'", ""),
	)
}
