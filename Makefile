CC=go
FMT=gofmt
NAME=bender
BASE_DIR=/script
OS=linux
ARCH=386
ARM=

default: fmt lint
	$(CC) build -o $(NAME) cmd/plugin.go

fmt:
	$(CC) fmt ./...

lint:
	golangci-lint run

check_fmt:
	test -z "$(gofmt -s -l $(find . -name '*.go' -type f -print) | tee /dev/stderr)"

check:
	$(CC) test -v ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) $(CC) build -o $(NAME)-$(OS)-$(ARCH) cmd/plugin.go

dep:
	$(CC) mod vendor
