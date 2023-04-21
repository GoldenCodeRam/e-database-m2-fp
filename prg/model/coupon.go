package model

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/goldencoderam/prg/database"
)

type Coupon struct {
	RequestProductId int
	discount         float32
	description      string
}

func (c Coupon) NewFakeRecord() (database.Insertable, error) {
	var requestProductId int

	database.ExecuteRowQuery(
		"SELECT id FROM request_product ORDER BY random() LIMIT 1",
		&requestProductId,
	)

	return Coupon{
		RequestProductId: requestProductId,
		discount:         rand.Float32(),
		description:      gofakeit.LoremIpsumSentence(5),
	}, nil
}

func (c Coupon) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO coupon (request_product_id, discount, description)
        VALUES (%d, %f, '%s')`,
		c.RequestProductId,
		c.discount,
		c.description,
	)
}
