.PHONY: build
build:
	go build -v ./cmd/webservice

.PHONY: test
test:
	go test -v -race -timeout 40s ./...


.DEFAULT_GOAL := build	

