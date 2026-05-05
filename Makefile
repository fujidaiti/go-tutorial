include .env

.PHONY: dev
dev:
	env $$(cat .env | xargs) go run cmd/web/*.go

.PHONY: db-migrate
db-migrate:
	env $$(cat .env | xargs) go run ./cmd/migrate

.PHONY: db-seed
db-seed:
	psql -h $(PGHOST) -p $(PGPORT) -U $(PGUSER) -d $(PGDATABASE) -f seed.sql

.PHONY: db-reset
db-reset:
	dropdb -h $(PGHOST) -p $(PGPORT) --if-exists $(PGDATABASE)
	createdb -h $(PGHOST) -p $(PGPORT) $(PGDATABASE)
	$(MAKE) db-migrate
	$(MAKE) db-seed
