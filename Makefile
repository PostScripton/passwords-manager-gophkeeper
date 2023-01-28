BIN := "./bin/cli"
APP_NAME="Passwords Manager GophKeeper"

GIT_HASH := $(shell git rev-parse HEAD)
LDFlAGS = -X 'main.buildVersion=v0.0.0' -X 'main.buildTime=$(shell date +'%Y-%m-%d %H:%M:%S')' -X 'main.buildCommit=$(GIT_HASH)'

build:
	go build -a -o $(BIN) -ldflags "$(LDFlAGS)" cmd/cli/main.go

run: build
	$(BIN) serve
