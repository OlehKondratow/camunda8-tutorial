.PHONY: build-example run-example

EXAMPLE := example-worker

build-example:
	go build -o $(EXAMPLE) ./cmd/example-worker

run-example: build-example
	set -a && test -f .env && . ./.env; set +a; \
	./$(EXAMPLE)
