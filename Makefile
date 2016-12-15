get-deps: vendor

vendor:
	glide install

test: overalls.coverprofile

overalls.coverprofile: get-deps
	overalls -project=github.com/twexler/gd-ddns-client

coverage.html: test
	go tool cover -html=overalls.coverprofile -o=coverage.html

clean:
	rm -rf vendor coverage.out coverage.html

all: clean test

.PHONY: clean get-deps test
