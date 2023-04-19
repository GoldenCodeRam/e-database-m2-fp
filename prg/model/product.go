package model

import (
	"fmt"

	"github.com/goldencoderam/prg/database"
)

type Product struct {
	Id string
}

func (p Product) NewFakeRecord() (database.Insertable, error) {
	return Product{
		Id: p.Id,
	}, nil
}

func (p Product) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO product (id)
        VALUES ('%s')`,
		p.Id,
	)
}
