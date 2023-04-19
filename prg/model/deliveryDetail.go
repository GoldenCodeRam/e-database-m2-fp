package model

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goldencoderam/prg/database"
)

type DeliveryDetail struct {
	RequestId   int
	Address     string
	Description string
}

func (d DeliveryDetail) NewFakeRecord() (database.Insertable, error) {
	var requestId int

	database.ExecuteRowQuery("SELECT id FROM request ORDER BY random() LIMIT 1", &requestId)

	return DeliveryDetail{
        RequestId: requestId,
        Address: gofakeit.Address().Address,
        Description: gofakeit.LoremIpsumSentence(5),
	}, nil
}

func (d DeliveryDetail) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO delivery_detail (request_id, address, description)
        VALUES (%d, '%s', '%s')`,
        d.RequestId,
        d.Address,
        d.Description,
	)
}
