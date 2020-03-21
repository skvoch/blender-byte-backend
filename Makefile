.PH0NY: build

build:
	go build -v ./cmd/back
	go build -v ./cmd/export
test:
	go test -v -race ./...

.DEFAULT_GOAL := build
