.PHONY: start-dev-mem start-dev-file start-dev-psql start-prod test test-cover lint install-deps
.SILENT:

start-dev-mem:
	clear; \
	ENV="dev" \
	go run ./cmd/shortener/main.go

start-dev-file:
	clear; \
	ENV="dev" \
	FILE_STORAGE_PATH="/tmp/short-url-db.json" \
	go run ./cmd/shortener/main.go

start-dev-psql:
	clear; \
	ENV="dev" \
	REST_WRITE_TIMEOUT="100000" \
	DATABASE_DSN="postgres://root:strongpass@localhost:3001/shortener?sslmode=disable" \
	go run ./cmd/shortener/main.go

start-prod:
	clear; \
	ENV="prod" \
	SERVER_ADDRESS="localhost:3000" \
	go run ./cmd/shortener/main.go


start-psql:
	@docker stop shortener_psql || true
	@docker run -d --rm \
		--name shortener_psql \
		-p 3001:5432 \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=strongpass \
		-e POSTGRES_DB=shortener \
		postgres:15.4-alpine3.17

test:
	go test ./... --cover

test-cover:
	go test --coverprofile=coverage.out ./... > /dev/null; \
    go tool cover -func=coverage.out | grep total | grep -oE '[0-9]+(\.[0-9]+)?%'

test-bench:
	go test ./internal/http-server/handler/... -bench=. -benchmem -memprofile=profiles/last.pprof

test-bench-show:
	go tool pprof -http=":9090" handler.test profiles/last.pprof

lint:
	@clear
	$(LOCAL_BIN_PATH)/golangci-lint run -c $(CONFIG) --path-prefix $(ROOT_PATH)

install-deps:
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/golang/mock/mockgen@v1.6.0
	@GOBIN=$(LOCAL_BIN_PATH)  go install golang.org/x/perf/cmd/benchstat@latest
	@go mod tidy

# -------

LOCAL_BIN_PATH := $(CURDIR)/bin
ROOT_PATH := $(realpath .)
CONFIG := $(ROOT_PATH)/.golangci.yaml