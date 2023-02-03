.PHONY: help

TAG=xdung24/repository
APPNAME=telegramchatbot
VERSION=$(shell cat VERSION)

help:  ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:  ## Build docker image
	@docker buildx build --no-cache --platform linux/amd64 --push  --tag $(TAG):$(APPNAME)-latest --tag $(TAG):$(APPNAME)-$(VERSION) -f Dockerfile .

test: ## Run tests
	@echo 'Running test...'
	@go test ./src
	@echo 'Done...'