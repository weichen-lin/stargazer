current_dir := $(shell pwd)
GITHUB_SHA ?= latest

run-pg:
	docker run --name finsurf \
	-d -p 5432:5432 -v finsurf-postgres:/var/lib/postgresql/data \
	-e POSTGRES_USER=finsurf -e POSTGRES_PASSWORD=finsurf \
	postgres:16.4

run-sqlc:
	docker run --rm -v $(current_dir):/src -w /src sqlc/sqlc:1.20.0 generate

create-migrate:
	docker run \
	-v $(current_dir)/db/migrate:/migrations \
	--network host migrate/migrate \
	-path=/migrations/ -database postgres://localhost:5432 \
	create -ext sql -dir migrations -seq 1

migrate-up:
	docker run \
	-v $(current_dir)/db/migrate:/migrations \
	--network host \
	migrate/migrate \
    -path=/migrations/ -database postgres://finsurf:finsurf@localhost:5432/stargazer?sslmode=disable up

migrate-down:
	docker run \
	-v $(current_dir)/db/migrate:/migrations \
	--network host \
	migrate/migrate \
    -path=/migrations/ -database postgres://finsurf:finsurf@localhost:5432/stargazer?sslmode=disable down -all

build-latest:
	docker build -t asia-east1-docker.pkg.dev/wei-dev/stargazer/backend:latest .

push-latest:
	docker push asia-east1-docker.pkg.dev/wei-dev/stargazer/backend:latest

replace-latest:
	yq -i ".spec.template.spec.containers[0].image = \"${GCP_PROJECT_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/${PROJECT_NAME}/backend:${IMAGE_TAG}\"" deploy/service.yaml && \
    gcloud beta run services replace deploy/service.yaml

apply-policy:
	gcloud run services add-iam-policy-binding finsurf-backend \
	--region=asia-east1 \
	--member="allUsers" \
	--role="roles/run.invoker"