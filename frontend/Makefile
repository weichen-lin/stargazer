
run-postgres:
	docker run --rm \
	-v stargazer-postgres:/var/lib/postgresql/data \
	-dp 5432:5432 \
	--name stargazer-postgres \
	-e POSTGRES_DB=local-dev \
	-e POSTGRES_USER=wcl \
	-e POSTGRES_PASSWORD=wcltestpwd \
	postgres:13-alpine \