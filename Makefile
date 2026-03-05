.PHONY: all
all:
	go build -o build

.PHONY: test
test:
	go test -v ./...
