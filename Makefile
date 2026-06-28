.DEFAULT_GOAL := build

clean:
	rm -rf bin

build: clean
	go build -o bin/mcp ./cmd/main.go