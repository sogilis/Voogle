include ../../../.env

run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 RABBITMQ_USER=$(RABBITMQ_USER) RABBITMQ_PWD=$(RABBITMQ_PWD) S3_HOST=http://localhost:9000 S3_AUTH_KEY=$(S3_AUTH_KEY) S3_AUTH_PWD=$(S3_AUTH_PWD) go run .

run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/encoder
build_image:
	docker build . -t voogle-encoder
