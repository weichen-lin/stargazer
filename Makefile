
run-postgres:
	docker run --rm \
	-v stargazer-postgres:/var/lib/postgresql/data \
	-dp 5432:5432 \
	--name stargazer-postgres \
	-e POSTGRES_DB=mydb \
	-e POSTGRES_USER=johndoe \
	-e POSTGRES_PASSWORD=randompassword \
	pgvector/pgvector:pg14

run-neo4j:
	docker run --rm \
	-v stargazer-neo4j:/data \
	-dp 7687:7687 \
	--name stargazer-neo4j \
	-e NEO4J_AUTH=neo4j/randompassword \
	neo4j:latest