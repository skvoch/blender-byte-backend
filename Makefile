.PH0NY: build

build:
	go build -v ./cmd/back
test:
	go test -v -race ./...

.DEFAULT_GOAL := build
