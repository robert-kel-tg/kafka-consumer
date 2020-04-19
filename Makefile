NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

BINARY_NAME=consumer
BINARY_SRC=./cmd/${BINARY_NAME}
GO_LINKER_FLAGS=-ldflags "-s"
DIR_OUT=$(CURDIR)/out

.PHONY: build install clean

build:
	@printf "$(OK_COLOR)==> Building binary$(NO_COLOR)\n"
	@go build -o ${DIR_OUT}/${BINARY_NAME} ${GO_LINKER_FLAGS} ${BINARY_SRC}

install:
	@printf "$(OK_COLOR)==> Installing binary$(NO_COLOR)\n"
	@go install -v ${BINARY_SRC}

clean:
	rm -rf vendor $(BUILD_DIR)
