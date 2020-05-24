NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

BINARY_NAME=consumer
BINARY_SRC=./cmd/${BINARY_NAME}
GO_LINKER_FLAGS=-ldflags "-s"
DIR_OUT=$(CURDIR)/out

PID_FILE := /tmp/$(BINARY_NAME).pid

OS:=linux
MIGRATE_TOOL_TAG=v4.11.0
MIGRATE_TOOL_URL:=https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_TOOL_TAG}/migrate.${OS}-amd64.tar.gz

.DEFAULT_GOAL := help

migrate:
	curl -L ${MIGRATE_TOOL_URL} | tar xvz
	# Move it to the target name
	mv migrate.${OS}-amd64 $@
	chmod +x $@

.PHONY: build
build: ## Builds the binary
	@printf "$(OK_COLOR)==> Building binary$(NO_COLOR)\n"
	@go build -o ${DIR_OUT}/${BINARY_NAME} ${GO_LINKER_FLAGS} ${BINARY_SRC}

.PHONY: install
install:
	@printf "$(OK_COLOR)==> Installing binary$(NO_COLOR)\n"
	@go install -v ${BINARY_SRC}

.PHONY: clean
clean: ## Cleans directory with the binary
	rm -rf vendor $(DIR_OUT)

.PHONY: dev-start
dev-start: ## Starts service for development
	@go build -o ${DIR_OUT}/${BINARY_NAME} ${BINARY_SRC} && ${DIR_OUT}/${BINARY_NAME}  & echo $$! > $(PID_FILE)
	@printf "$(OK_COLOR)==> Starting service of: $(NO_COLOR)PID $$(cat $(PID_FILE))\n"

.PHONY: dev-stop
dev-stop: ## Stops service and kills PID tree
	@printf "$(OK_COLOR)==> Killing service of: $(NO_COLOR)PID $$(cat $(PID_FILE))\n"
	@-kill `pstree -p \`cat $(PID_FILE)\` | tr "\n" " " |sed "s/[^0-9]/ /g" |sed "s/\s\s*/ /g"`

.PHONY: dev-restart
dev-restart: dev-stop dev-start ## Restarts, stops and starts again the service

.PHONY: dev-run
dev-run: dev-start ## Runs the service with reload (once after files got changed)
	@fswatch -or --event=Updated . | xargs -n1 -I {} make dev-restart

.PHONY: migrate-up
migrate-up: migrate ## Runs DB migrations
	./migrate -database=${DB_DSN} -path migrations up

.PHONY: test-unit
test-unit:  ## Run unit tests
	go test -v -race -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: help
help: ## List of all available targets (Note: this is default target)
	@echo List of all available target commands:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
