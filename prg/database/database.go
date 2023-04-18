package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func doExecWithConnection(do func(tx pgx.Tx) error) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal(err, "unable to connect to database")
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatal(err, "unable to start transaction")
	}

	tx.Begin(context.Background())
	err = do(tx)
	if err != nil {
		tx.Rollback(context.Background())
	}
	tx.Commit(context.Background())

	conn.Close()

	return err
}

func ExecuteFakeInsert(
	record FakeRecord,
	amount int,
) error {
	err := doExecWithConnection(func(tx pgx.Tx) error {
		for i := 0; i < amount; i++ {
			log.Println(fmt.Sprintf("Fake insert %d of %d", i, amount))

			val, err := record.NewFakeRecord()
			if err != nil {
				log.Fatal("unable to create fake record")
			}
			_, err = tx.Exec(
				context.Background(),
				val.GenerateInsertSql(),
			)

			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
