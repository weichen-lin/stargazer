run-kafka:
	docker-compose -f kafka.docker-compose.yml up -d

stop-kafka:
	docker-compose -f kafka.docker-compose.yml down

format:
	find $(PWD) -name "*.go" -exec gofmt -w {} \;

migrate-up:
	migrate \
	-path db/migration \
	-database "postgresql://johndoe:randompassword@localhost:5432/mydb?sslmode=disable" \
	-verbose \
	up

migrate-down:
	migrate \
	-path db/migration \
	-database "postgresql://johndoe:randompassword@localhost:5432/mydb?sslmode=disable" \
	-verbose \
	down

# ./kafka-console-producer.sh --bootstrap-server your_kafka_broker:9092 --topic your_topic_name
