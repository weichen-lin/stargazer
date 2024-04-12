location = asia-east1
project_id = stargazer-414509
repository = stargazer-images

run-kafka:
	docker-compose -f kafka.docker-compose.yml up

stop-kafka:
	docker-compose -f kafka.docker-compose.yml down

format:
	find $(PWD) -name "*.go" -exec gofmt -w {} \;

# ./kafka-console-producer.sh --bootstrap-server your_kafka_broker:9092 --topic your_topic_name

build-kafka:
	docker build -t $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-producer:latest .

push-kafka: 
	docker push $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-producer:latest