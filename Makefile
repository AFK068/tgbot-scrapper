COVERAGE_FILE ?= coverage.out

.PHONY: run_bot
run_bot:
	@echo "Выполняется go run для таргета bot"
	@go run ./cmd/bot/main.go

.PHONY: run_scrapper
run_scrapper:
	@echo "Выполняется go run для таргета scrapper"
	@go run ./cmd/scrapper/main.go

.PHONY: build
build: build_bot build_scrapper

.PHONY: build_bot
build_bot:
	@echo "Выполняется go build для таргета bot"
	@mkdir -p .bin
	@go build -o ./bin/bot ./cmd/bot

.PHONY: build_scrapper
build_scrapper:
	@echo "Выполняется go build для таргета scrapper"
	@mkdir -p .bin
	@go build -o ./bin/scrapper ./cmd/scrapper


## test: run all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/es-debug/backend-academy-2024-go-template/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'

.PHONY: lint
lint: lint-golang lint-proto

.PHONY: lint-golang
lint-golang:
	@if ! command -v 'golangci-lint' &> /dev/null; then \
  		echo "Please install golangci-lint!"; exit 1; \
  	fi;
	@golangci-lint -v run --fix ./...

.PHONY: lint-proto
lint-proto:
	@if ! command -v 'easyp' &> /dev/null; then \
  		echo "Please install easyp!"; exit 1; \
	fi;
	@easyp lint

.PHONY: generate
generate: generate_proto generate_openapi

.PHONY: generate_proto
generate_proto:
	@if ! command -v 'easyp' &> /dev/null; then \
		echo "Please install easyp!"; exit 1; \
	fi;
	@easyp generate

.PHONY: generate_openapi
generate_openapi:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;
	@mkdir -p internal/api/openapi/scrapper/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		api/openapi/v1/scrapper-api.yaml > internal/api/openapi/scrapper/v1/scrapper-api.gen.go
	@mkdir -p internal/api/openapi/bot/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		api/openapi/v1/bot-api.yaml > internal/api/openapi/bot/v1/bot-api.gen.go

.PHONY: clean
clean:
	@rm -rf./bin