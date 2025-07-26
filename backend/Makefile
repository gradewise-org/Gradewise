APP_SERVER=cmd/server
APP_WORKER=cmd/worker
BIN_DIR=bin

.PHONY: all server worker clean

all: server worker

internal/api/gen.go: api/api.yaml api/oapi-codegen.yaml
	@echo "Generating API code from OpenAPI spec..."
	go generate

generate: internal/api/gen.go

server: generate
	@echo "Building server..."
	time go build -o $(BIN_DIR)/server $(APP_SERVER)/main.go
	@ls -lh $(BIN_DIR)/server

worker: generate
	@echo "Building worker..."
	time go build -o $(BIN_DIR)/worker $(APP_WORKER)/main.go
	@ls -lh $(BIN_DIR)/worker

clean:
	rm  $(BIN_DIR)/*
	rm internal/api/gen.go