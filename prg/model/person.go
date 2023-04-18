package model

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goldencoderam/prg/database"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Address   string
}

func (p Person) NewFakeRecord() (database.Insertable, error) {
    cleanPassword := gofakeit.Word()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cleanPassword), 10)

	if err != nil {
		return nil, err
	}

	return Person{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  string(hashedPassword),
		Address:   gofakeit.Address().Address,
	}, nil
}

func (p Person) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO person (first_name, last_name, email, password, address)
        VALUES ('%s', '%s', '%s', '%s', '%s')`,
		p.FirstName,
		p.LastName,
		p.Email,
		p.Password,
		p.Address,
	)
}
