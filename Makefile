NAME := lonely

all: build-linux-arm

build-linux-arm:
	@$(MAKE) build GOOS=linux GOARCH=arm

.PHONY: test build

build:
	go build -o build/$(GOOS)_$(GOARCH)/$(NAME)

test:
	go test

clean:
	rm build/*
