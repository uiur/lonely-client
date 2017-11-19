.PHONY: test build clean installdeps

NAME := lonely

GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(THIS_GOOS)
GOARCH=$(THIS_GOARCH)

all: build-linux-arm

build-linux-arm:
	@$(MAKE) build GOOS=linux GOARCH=arm

build:
	go build -o build/$(GOOS)_$(GOARCH)/$(NAME)

test:
	go test -v ./...

clean:
	rm -rf build/*

installdeps:
	go get -t -d -v ./...
