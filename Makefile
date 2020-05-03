NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

BINARY_NAME=consumer
BINARY_SRC=./cmd/${BINARY_NAME}
GO_LINKER_FLAGS=-ldflags "-s"
DIR_OUT=$(CURDIR)/out

PID_FILE := /tmp/$(BINARY_NAME).pid

.PHONY: build install clean

build:
	@printf "$(OK_COLOR)==> Building binary$(NO_COLOR)\n"
	@go build -o ${DIR_OUT}/${BINARY_NAME} ${GO_LINKER_FLAGS} ${BINARY_SRC}

install:
	@printf "$(OK_COLOR)==> Installing binary$(NO_COLOR)\n"
	@go install -v ${BINARY_SRC}

clean:
	rm -rf vendor $(BUILD_DIR)

dev-start:
	@go build -o ${DIR_OUT}/${BINARY_NAME} ${BINARY_SRC} && ${DIR_OUT}/${BINARY_NAME}  & echo $$! > $(PID_FILE)
	@printf "$(OK_COLOR)==> Starting service of: $(NO_COLOR)PID $$(cat $(PID_FILE))\n"

dev-stop:
	@printf "$(OK_COLOR)==> Killing service of: $(NO_COLOR)PID $$(cat $(PID_FILE))\n"
	@-kill `pstree -p \`cat $(PID_FILE)\` | tr "\n" " " |sed "s/[^0-9]/ /g" |sed "s/\s\s*/ /g"`

dev-restart: dev-stop dev-start

dev-run: dev-start
	@fswatch -or --event=Updated . | xargs -n1 -I {} make dev-restart
