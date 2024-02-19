location = asia-east1
project_id = stargazer-414509
repository = stargazer-images

auth-docker:
	cat credentials.json | docker login -u _json_key_base64 --password-stdin \
	https://asia-east1-docker.pkg.dev

build-base:
	docker build -t $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-transformer:base -f base.Dockerfile .

push-base: 
	docker push $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-transformer:base

build-transformer:
	docker build \
    $(echo $(cat secrets.env | xargs -I {} echo "--build-arg {}")) \
    -t $(location)-docker.pkg.dev/$(project_id)/$(repository)/stargazer-transformer:latest \
    .
