EXECUTABLE=zrun
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: build test

test:
	go test ./...

build: windows linux darwin
	@echo version: $(VERSION)

windows: $(WINDOWS)

linux: $(LINUX)

darwin: $(DARWIN)

$(WINDOWS):
	@go env GOOS=windows
	@go env GOARCH=amd64
	@go build -v -o bin\$(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" main.go

$(LINUX):
	@go env GOARCH=amd64
	@go env GOOS=linux
	@go build -v -o bin\$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" main.go

$(DARWIN):
	@go env GOOS=darwin
	@go env GOARCH=amd64
	@go build -v -o bin\$(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" main.go

clean:
	@del /f bin\$(WINDOWS)
	@del /f bin\$(LINUX)
	@del /f bin\$(DARWIN)

run-ubuntu: clean build
	@echo "Running Ubuntu"
	@docker build -t zrun-ubuntu -f .\scripts\docker\ubuntu\Dockerfile .
	@docker-compose -f .\scripts\docker\ubuntu\docker-compose.yaml down -v
	@docker-compose -f .\scripts\docker\ubuntu\docker-compose.yaml up -d
	@docker exec -it zrun-ubuntu /bin/bash -c "./zrun about"
exec-ubuntu:
	@docker exec -it zrun-ubuntu /bin/bash
