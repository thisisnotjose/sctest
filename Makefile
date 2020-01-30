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
	github.com/thisisnot/sctest/cmd/ \
	github.com/thisisnot/sctest/internal/processors/ \
	github.com/thisisnot/sctest/internal/users/ \
	github.com/thisisnot/sctest/internal/handlers/ \
	-coverprofile=/tmp/code-cover