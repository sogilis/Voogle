include ../../../.env

run:
	go run .
run-dev:
	DEV_MODE=true S3_HOST=http://localhost:9000 S3_AUTH_KEY=$(S3_AUTH_KEY) S3_AUTH_PWD=$(S3_AUTH_PWD) CONSUL_HOST=$(CONSUL_HOST) go run .

run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/flip-server-transformer
build_image:
	docker build . -t voogle-flip-server-transformer
