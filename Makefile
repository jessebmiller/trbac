# Dependency managemnt lifecycle

LOCAL_BUILD_NAME = $(shell basename $$(git rev-parse --show-toplevel)):local


.PHONY: ensure
ensure: build
	docker run -v $(shell pwd):$(shell docker run $(LOCAL_BUILD_NAME) pwd) $(LOCAL_BUILD_NAME) dep ensure


.PHONY: build
build:
	docker build -t $(LOCAL_BUILD_NAME) .
