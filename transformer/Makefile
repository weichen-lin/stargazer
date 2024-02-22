location = asia-east1
project_id = stargazer-414509
repository = stargazer-images

auth-docker:
	cat credentials.json | docker login -u _json_key_base64 --password-stdin \
	https://asia-east1-docker.pkg.dev

build-transformer:
	docker build -t $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-transformer:latest .

push-transformer:
	docker push $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-transformer:latest


# unset AUTHENTICATION_TOKEN DATABASE_URL OPENAI_API_KEY VSCODE_ENV_REPLACE