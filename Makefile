.PHONY: all

export_name = "redis-protocol-test"

build:
	go build -o $(export_name) main.go

all: build

Run: build
	./$(export_name)