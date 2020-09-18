GOPATH=
GO111MODULE=on
GOPROXY=https://goproxy.cn,direct

GO = CGO_ENABLED=0 go
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_BRANCH := $(shell git symbolic-ref --short -q HEAD)
GIT_COMMIT_HASH := $(shell git rev-parse --verify HEAD)
GO_FLAGS := -v -ldflags="-X 'github.com/voidint/gbb/build.Date=$(BUILD_DATE)' -X 'github.com/voidint/gbb/build.Commit=$(GIT_COMMIT_HASH)' -X 'github.com/voidint/gbb/build.Branch=$(GIT_BRANCH)'"

all: install test

build:
	@echo "GO111MODULE=$(GO111MODULE)"
	@echo "GOPROXY=$(GOPROXY)"
	$(GO) build $(GO_FLAGS)

install: build
	$(GO) install $(GO_FLAGS)

test:
	$(GO) test -v ./...

clean:
	$(GO) clean -x

.PHONY: all build install test clean