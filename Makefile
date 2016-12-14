get-deps: vendor

vendor:
	glide install

test: get-deps coverage.out

coverage.out:
	go test -coverprofile coverage.out

coverage.html: test
	go tool cover -html=coverage.out -o=coverage.html

clean:
	rm -rf vendor coverage.out coverage.html

all: clean test

.PHONY: clean get-deps test
