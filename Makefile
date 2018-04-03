# Dependency managemnt lifecycle

BUILD_NAME = $(shell basename $$(git rev-parse --show-toplevel)):$RANDOM


.PHONY: ensure
ensure: build
	docker run -v $(shell pwd):$(shell docker run $(BUILD_NAME) pwd) $(BUILD_NAME) dep ensure


.PHONY: build
build:
	docker build -t $(BUILD_NAME) .
