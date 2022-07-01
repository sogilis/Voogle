include ../../../.env

run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 RABBITMQ_USER=$(RABBITMQ_USER) RABBITMQ_PWD=$(RABBITMQ_PWD) S3_HOST=http://localhost:9000 S3_AUTH_KEY=$(S3_AUTH_KEY) S3_AUTH_PWD=$(S3_AUTH_PWD) USER_AUTH=$(USER_AUTH) PWD_AUTH=$(PWD_AUTH) MARIADB_HOST=localhost MARIADB_PORT=$(MARIADB_PORT) MARIADB_USER=$(MARIADB_USER) MARIADB_PASSWORD=$(MARIADB_PASSWORD) MARIADB_DATABASE=$(MARIADB_DATABASE) CONSUL_HOST=localhost:8500 go run .
run-dev-remote:
	DEV_MODE=true USER_AUTH=$(USER_AUTH) PWD_AUTH=$(PWD_AUTH) go run .

build:
	go build -o build/api
build_image:
	docker build . -t voogle-api

generate-docs:
	swag init -g main.go
