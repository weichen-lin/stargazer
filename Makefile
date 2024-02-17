
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
	-p7474:7474 -p7687:7687 \
	-d \
	--name stargazer-neo4j \
	-e NEO4J_AUTH=neo4j/randompassword \
	neo4j:latest

migrate-up:
	migrate \
	-path migrate \
	-database "postgresql://johndoe:randompassword@localhost:5432/mydb?sslmode=disable" \
	-verbose \
	up

migrate-down:
	migrate \
	-path migrate \
	-database "postgresql://johndoe:randompassword@localhost:5432/mydb?sslmode=disable" \
	-verbose \
	down