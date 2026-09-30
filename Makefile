.PHONY: cli-install test

cli-install:
	@go install

test:
	@go test -short -cover ./...
