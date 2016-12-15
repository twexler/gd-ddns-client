BIN_FILES=$(shell find . -name vendor -prune -o -not -name '*_test.go' -name '*.go' -print)
SRC_FILES=$(concat $(shell find . -name vendor -prune -o -name '*.go' -print) $(BIN_FILES))
get-deps: vendor

vendor:
	glide install

test: overalls.coverprofile

coverage: coverage.html

overalls.coverprofile: get-deps $(SRC_FILES)
	overalls -project=github.com/twexler/gd-ddns-client

coverage.html: overalls.coverprofile
	go tool cover -html=overalls.coverprofile -o=coverage.html

clean:
	rm -rf vendor coverage.html gd-ddns-client

gd-ddns-client: $(BIN_FILES)
	go build cmd/gd-ddns-client.go

all: clean test

.PHONY: clean coverage docker-image get-deps test
