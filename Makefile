OS_NAME := $(shell uname)
ifeq ($(OS_NAME), Darwin)
OPEN := open
else
OPEN := xdg-open
endif

build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/ide-config-sync.exe ./main.go

	@GOOS=linux GOARCH=amd64 go build -o ./bin/ide-config-sync-darwin-amd64 ./main.go
