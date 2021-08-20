PROJECT=ipfs_repl

all: build

build:
	go build -o bin/main main.go

start: 
	@./bin/main
