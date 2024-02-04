run-kafka:
	docker-compose -f kafka.docker-compose.yml up -d

stop-kafka:
	docker-compose -f kafka.docker-compose.yml down

format:
	find $(PWD) -name "*.go" -exec gofmt -w {} \;