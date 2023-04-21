#!make
include .env
export

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path ${MIGRATIONS_PATH} up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path ${MIGRATIONS_PATH} down

migrate-force:
	migrate -database ${POSTGRESQL_URL} -path ${MIGRATIONS_PATH} force $(VERSION)

migrate-create:
	migrate create -ext sql -dir ${MIGRATIONS_PATH} -seq $(NAME)

ferret-up:
	docker run --rm --name ferret -d -p 9222:9222 montferret/chromium

ferret-down:
	docker stop ferret -t 0
