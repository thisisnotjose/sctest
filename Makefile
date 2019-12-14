.PHONY: build clean

start:
	./bin/main

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/main src/main.go

clean:
	rm -rf ./bin ./vendor

test:
	GO111MODULE=on go test -timeout 30s \
	github.com/thisisnotjose/sctest/main \
	-coverprofile=/tmp/code-cover