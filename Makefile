.PHONY: cli-install format goimports-install test
.PHONY: migration-create migration-install migration-up

MIGRATIONS_DIR := db-migrations

cli-install:
	@go install

format:
	@goimports -w .

go-imports-install:
	@go install golang.org/x/tools/cmd/goimports@latest

migration-create:
	# usage:
	# make migration-create name=add_users_table
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -format "20060102150405" $(name)

migration-install:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migration-up:
	@migrate -source file://db-migrations -database "postgres://postgres:password@localhost/jira?sslmode=disable" up

test:
	@go test -cover -short ./...
