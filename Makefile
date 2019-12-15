.PHONY: build clean

start: build
	./bin/main

build:
	export GO111MODULE=on
	env go build -ldflags="-s -w" -o bin/main cmd/main.go

clean:
	rm -rf ./bin ./vendor

test:
	GO111MODULE=on go test -timeout 30s \
	github.com/thisisnotjose/sctest/cmd/ \
	github.com/thisisnotjose/sctest/internal/processors/ \
	github.com/thisisnotjose/sctest/internal/handlers/ \
	-coverprofile=/tmp/code-cover