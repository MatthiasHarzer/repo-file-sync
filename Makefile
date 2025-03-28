OS_NAME := $(shell uname)
ifeq ($(OS_NAME), Darwin)
OPEN := open
else
OPEN := xdg-open
endif

BUILD_VERSION ?= "unknown"

build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/repo-file-sync.exe -ldflags "-X main.version=$(BUILD_VERSION)" ./main.go

	@GOOS=linux GOARCH=amd64 go build -o ./bin/repo-file-sync -ldflags "-X main.version=$(BUILD_VERSION)" ./main.go

.PHONY: build
