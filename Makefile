
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

migrate-down-one:
	migrate \
	-path migrate \
	-database "postgresql://johndoe:randompassword@localhost:5432/mydb?sslmode=disable" \
	-verbose \
	down 1

service-up:
	docker-compose -f service.docker-compose.yml up -d

service-down:
	docker-compose -f service.docker-compose.yml down

producer-up:
	docker-compose -f kafka.docker-compose.yml up

producer-down:
	docker-compose -f kafka.docker-compose.yml down

clean-env:
	unset AUTHENTICATION_TOKEN DATABASE_URL OPENAI_API_KEY VSCODE_ENV_REPLACE GITHUB_CLIENT_ID GITHUB_CLIENT_SECRET NEXTAUTH_URL NEO4J_URL NEO4J_PASSWORD