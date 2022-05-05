USER_AUTH = $(shell grep -oP "USER_AUTH=\K.*" ../../../.env)
PWD_AUTH = $(shell grep -oP "PWD_AUTH=\K.*" ../../../.env)
RABBITMQ_USER = $(shell grep -oP "RABBITMQ_USER=\K.*" ../../../.env)
RABBITMQ_PWD = $(shell grep -oP "RABBITMQ_PWD=\K.*" ../../../.env)
S3_AUTH_KEY = $(shell grep -oP "S3_AUTH_KEY=\K.*" ../../../.env)
S3_AUTH_PWD = $(shell grep -oP "S3_AUTH_PWD=\K.*" ../../../.env)
MARIADB_DATABASE = $(shell grep -oP "MARIADB_DATABASE=\K.*" ../../../.env)
MARIADB_USER = $(shell grep -oP "MARIADB_USER=\K.*" ../../../.env)
MARIADB_PASSWORD = $(shell grep -oP "MARIADB_PASSWORD=\K.*" ../../../.env)
MARIADB_PORT = $(shell grep -oP "MARIADB_PORT=\K.*" ../../../.env)

run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 RABBITMQ_USER=$(RABBITMQ_USER) RABBITMQ_PWD=$(RABBITMQ_PWD) S3_HOST=http://localhost:9000 S3_AUTH_KEY=$(S3_AUTH_KEY) S3_AUTH_PWD=$(S3_AUTH_PWD) USER_AUTH=$(USER_AUTH) PWD_AUTH=$(PWD_AUTH) MARIADB_HOST=localhost MARIADB_PORT=$(MARIADB_PORT) MARIADB_USER=$(MARIADB_USER) MARIADB_PASSWORD=$(MARIADB_PASSWORD) MARIADB_DATABASE=$(MARIADB_DATABASE) go run .
run-dev-remote:
	DEV_MODE=true USER_AUTH=$(USER_AUTH) PWD_AUTH=$(PWD_AUTH) go run .

build:
	go build -o build/api
build_image:
	docker build . -t voogle-api

generate-docs:
	swag init -g main.go
