GIT_SHA := $(shell git rev-parse --short HEAD)
NAMESPACE := weichen-lin
REPO_NAME := stargazer
TAG := $(GIT_SHA)
IMAGE_BASE := ghcr.io/$(NAMESPACE)/$(REPO_NAME)

format:
	find $(PWD) -name "*.go" -exec gofmt -w {} \;

build-dev:
	docker build -t $(IMAGE_BASE)/kafka-service:latest .

push-dev:
	docker push $(IMAGE_BASE)/kafka-service:latest