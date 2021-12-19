PROJECT=ipfs_repl

all: build

build:
	go build -o bin/main ./...

start: 
	@./bin/main
