.PHONY: build run

build:
	go build

run: build
	./gocmdtimeout
