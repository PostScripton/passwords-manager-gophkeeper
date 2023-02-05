BIN := "./bin/cli"
APP_NAME="Passwords Manager GophKeeper"

GIT_HASH := $(shell git rev-parse HEAD)
LDFlAGS = -X 'main.buildVersion=v0.0.0' -X 'main.buildTime=$(shell date +'%Y-%m-%d %H:%M:%S')' -X 'main.buildCommit=$(GIT_HASH)'

build:
	go build -a -o $(BIN) -ldflags "$(LDFlAGS)" cmd/cli/main.go

run: build
	$(BIN) serve

#example: make proto name=user
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/$(name).proto
