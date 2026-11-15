.PHONY: all build test clean run-debate run-server web-build

BINARY_NAME=consensai.exe
BIN_DIR=bin

all: build

web-build:
	cd web && npm install && npm run build

build: web-build
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/consensai

test:
	go test -v ./...

run-debate:
	go run ./cmd/consensai debate --dossier examples/dossiers/commercial_loan_500k.json

run-server:
	go run ./cmd/consensai serve --port 6060

clean:
	rm -rf $(BIN_DIR)
	rm -rf web/dist
