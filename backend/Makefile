APP_SERVER=cmd/server
APP_WORKER=cmd/worker
BIN_DIR=bin

.PHONY: all server worker clean

all: server worker

server:
	@echo "Building server..."
	time go build -o $(BIN_DIR)/server $(APP_SERVER)/main.go
	@ls -lh $(BIN_DIR)/server
	@echo ""

worker:
	@echo "Building worker..."
	time go build -o $(BIN_DIR)/worker $(APP_WORKER)/main.go
	@ls -lh $(BIN_DIR)/worker
	@echo ""

clean:
	rm -rf $(BIN_DIR)/*