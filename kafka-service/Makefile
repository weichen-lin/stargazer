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

run-neo4j:
	docker run --rm -d --name stargazer_neo4j \
		-p 7474:7474 -p 7687:7687 \
		-e NEO4J_AUTH=neo4j/password \
		-v neo4j-data:/data \
		neo4j:5.18.0