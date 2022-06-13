export CGO_ENABLED=0

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install ./...

.PHONY: build
build:
	go build ./...
