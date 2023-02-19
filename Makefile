CLI_BIN := "./bin/cli"
SERVER_BIN := "./bin/server"
APP_NAME="Passwords Manager GophKeeper"

GIT_HASH := $(shell git rev-parse HEAD)
LDFlAGS = -X 'main.buildVersion=v0.0.0' -X 'main.buildTime=$(shell date +'%Y-%m-%d %H:%M:%S')' -X 'main.buildCommit=$(GIT_HASH)'

build-cli:
	go build -a -o $(CLI_BIN) -ldflags "$(LDFlAGS)" cmd/cli/main.go

build-server:
	go build -v -o $(SERVER_BIN) cmd/server/main.go

#example: make proto name=user
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/$(name).proto

cert-gen:
	cd ./cert; sh gen.sh
