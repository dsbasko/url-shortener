.PHONY: all
.SILENT:

start-http-dev-mem:
	clear; \
	ENV="dev" \
	CONFIG="config/shortener.json" \
	go run ./cmd/shortener

start-http-prod:
	clear; \
	ENV="prod" \
	CONFIG="config/shortener.json" \
	ENABLE_HTTPS="true" \
	SERVER_ADDRESS="localhost:3000" \
	go run ./cmd/shortener

start-http-dev-file:
	clear; \
	ENV="dev" \
	CONFIG="config/shortener.json" \
	FILE_STORAGE_PATH="/tmp/short-url-db.json" \
	go run ./cmd/shortener

start-http-dev-psql:
	clear; \
	ENV="dev" \
	CONFIG="config/shortener.json" \
	REST_WRITE_TIMEOUT="100000" \
	DATABASE_DSN="postgres://$(PSQL_USER):$(PSQL_PASS)@localhost:$(PSQL_PORT)/$(PSQL_DB)?sslmode=disable" \
	go run ./cmd/shortener

start-grpc-dev-mem:
	clear; \
	ENV="dev" \
	CONTROLLER="grpc" \
	CONFIG="config/shortener.json" \
	go run ./cmd/shortener

start-grpc-dev-file:
	clear; \
	ENV="dev" \
	CONTROLLER="grpc" \
	CONFIG="config/shortener.json" \
	FILE_STORAGE_PATH="/tmp/short-url-db.json" \
	go run ./cmd/shortener

start-grpc-dev-psql:
	clear; \
	ENV="dev" \
	CONTROLLER="grpc" \
	CONFIG="config/shortener.json" \
	REST_WRITE_TIMEOUT="100000" \
	DATABASE_DSN="postgres://$(PSQL_USER):$(PSQL_PASS)@localhost:$(PSQL_PORT)/$(PSQL_DB)?sslmode=disable" \
	go run ./cmd/shortener

start-grpc-prod:
	clear; \
	ENV="prod" \
	CONTROLLER="grpc" \
	CONFIG="config/shortener.json" \
	ENABLE_HTTPS="true" \
	SERVER_ADDRESS="localhost:3000" \
	go run ./cmd/shortener


start-psql:
	@docker stop shortener_psql || true
	@docker run -d --rm \
		--name shortener_psql \
		-p $(PSQL_PORT):5432 \
		-e POSTGRES_USER=$(PSQL_USER) \
		-e POSTGRES_PASSWORD=$(PSQL_PASS) \
		-e POSTGRES_DB=$(PSQL_DB) \
		postgres:15.4-alpine3.17

lint:
	@clear
	@$(LOCAL_BIN_PATH)/golangci-lint run -c $(CONFIG) --path-prefix $(ROOT_PATH)
	@go run $(CURDIR)/cmd/staticlint ./...

install-deps:
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	@GOBIN=$(LOCAL_BIN_PATH) go install go.uber.org/mock/mockgen@latest
	@GOBIN=$(LOCAL_BIN_PATH) go install golang.org/x/perf/cmd/benchstat@latest
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/nikolaydubina/go-cover-treemap@latest
	@GOBIN=$(LOCAL_BIN_PATH) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=$(LOCAL_BIN_PATH) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go mod tidy

build-staticlint:
	@go build -o $(LOCAL_BIN_PATH)/staticlint ./cmd/staticlint/main.go

test-cover:
	@go test --coverprofile=coverage.out ./... > /dev/null
	@go tool cover -func=coverage.out | grep total | grep -oE '[0-9]+(\.[0-9]+)?%'

test-cover-svg:
	@clear;
	@go test --coverprofile=coverage.out ./... > /dev/null;
	@$(LOCAL_BIN_PATH)/go-cover-treemap -coverprofile coverage.out > coverage.svg
	@xdg-open ./coverage.svg

test-cover-html:
	@clear;
	@go test --coverprofile=coverage.out ./... > /dev/null;
	@go tool cover -html="coverage.out"

test-bench:
	@go test ./internal/controller/rest/handlers -bench=. -benchmem -cpuprofile=profiles/cpu-last.pprof -memprofile=profiles/mem-last.pprof

test-bench-show-cpu:
	@go test ./internal/controller/rest/handlers -bench=. -benchmem -cpuprofile=profiles/cpu-last.pprof
	@go tool pprof -http=":9090" handler.test profiles/cpu-last.pprof

test-bench-show-mem:
	@go test ./internal/controller/rest/handlers -bench=. -benchmem -memprofile=profiles/mem-last.pprof
	@go tool pprof -http=":9090" handler.test profiles/mem-last.pprof

auto-tests:
	@clear
	@go build -o $(LOCAL_BIN_PATH)/shortener $(ROOT_PATH)/cmd/shortener/*.go
	@$(ROOT_PATH)/bin/shortenertest -test.v -test.run=$(TEST_RUN_ITERATION) "$(AUTO_TEST_CONFIG)"
	@printf "\n%.0s" $(seq 1 5)
	@go vet -vettool $(ROOT_PATH)/bin/statictest ./...

gen-cert:
	@go run ./cmd/cert-gen

gen-proto:
	@mkdir -p api/proto
	@protoc \
		--go_out=api --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN_PATH)/protoc-gen-go \
		--go-grpc_out=api --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN_PATH)/protoc-gen-go-grpc \
		proto/url_shortener_v1.proto

# -------
# Configs
# -------

LOCAL_BIN_PATH := $(CURDIR)/bin
ROOT_PATH := $(realpath .)
CONFIG := $(ROOT_PATH)/.golangci.yaml
SERVER_PORT := $(shell $($ROOT_PATH) bin/random unused-port)
TEMP_FILE := $(shell $(ROOT_PATH)/bin/random tempfile)

PSQL_USER := root
PSQL_PASS := strongpass
PSQL_DB := shortener
PSQL_PORT := 3091

AUTO_TEST_CONFIG := "-binary-path=$(ROOT_PATH)/bin/shortener \
	-server-port=$(SERVER_PORT) \
	-source-path=$(ROOT_PATH) \
	-file-storage-path=$(TEMP_FILE) \
	-database-dsn='postgresql://$(PSQL_USER):$(PSQL_PASS)@localhost:$(PSQL_PORT)/$(PSQL_DB)?sslmode=disable'"

TEST_RUN_ITERATION := ^TestIteration\([1-9]\|[1][0-9]\|2[0-5]\)$$
ifdef only
	TEST_RUN_ITERATION := ^TestIteration$(only)$$
endif