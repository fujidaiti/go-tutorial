.PHONY: db-migrate

db-migrate:
	env $(cat .env | xargs) go run ./cmd/migrate
