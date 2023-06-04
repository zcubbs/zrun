EXECUTABLE=zrun
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)
VAGRANT_DIR=scripts/vagrant/ubuntu
CURRENT_DIR=$(shell pwd)

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

d: build-docker-ubuntu run-docker-ubuntu
db: build-docker-ubuntu
dk: kill-docker-ubuntu
dx: exec-ubuntu
dd: mount-on-docker

build-docker-ubuntu:
	@docker build -t zrun-ubuntu -f .\scripts\docker\ubuntu\Dockerfile .

run-docker-ubuntu: kill-docker-ubuntu clean build
	@echo "Running Ubuntu"
	@docker-compose -f .\scripts\docker\ubuntu\docker-compose.yaml up -d
	@docker exec -it zrun-ubuntu /bin/bash -c "./zrun about"

kill-docker-ubuntu:
	@docker-compose -f .\scripts\docker\ubuntu\docker-compose.yaml down -v

mount-on-docker:
	@docker rm -f zrun-live
	@echo "${CURDIR}"
	@docker run -d -v "${CURDIR}":/app -w /app --name zrun-live zrun-ubuntu tail -f /dev/null
	@docker exec -it zrun-live /bin/bash

exec-ubuntu:
	@docker exec -it zrun-ubuntu /bin/bash

lint:
	@golangci-lint run ./... -v --timeout 5m

scan:
	@gosec -conf .gosec.config.json -fmt=json -out=gosec-results.json "./..."

v: vagrant-ubuntu-reload
vssh: vagrant-ubuntu-ssh
vfssh: vagrant-ubuntu vagrant-ubuntu-ssh
vk: kill-vagrant-ubuntu

vagrant-ubuntu: kill-vagrant-ubuntu vagrant-ubuntu-reload

vagrant-ubuntu-reload: clean build
	@echo "Running Ubuntu"
	cd $(VAGRANT_DIR) && vagrant up --provision

vagrant-ubuntu-ssh:
	cd $(VAGRANT_DIR) && vagrant ssh

kill-vagrant-ubuntu:
	cd $(VAGRANT_DIR) && vagrant destroy -f
