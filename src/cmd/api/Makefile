USER_AUTH = $(shell cat ../../../.env | grep USER_AUTH | awk -F'=' '{print $2}' | tr -d '\n')
PWD_AUTH = $(shell cat ../../../.env | grep PWD_AUTH | awk -F'=' '{print $2}' | tr -d '\n')
RABBITMQ_USER = $(shell cat ../../../.env | grep RABBITMQ_USER | awk -F'=' '{print $2}' | tr -d '\n')
RABBITMQ_PWD = $(shell cat ../../../.env | grep RABBITMQ_PWD | awk -F'=' '{print $2}' | tr -d '\n')
S3_AUTH_KEY = $(shell cat ../../../.env | grep S3_AUTH_KEY | awk -F'=' '{print $2}' | tr -d '\n')
S3_AUTH_PWD = $(shell cat ../../../.env | grep S3_AUTH_PWD | awk -F'=' '{print $2}' | tr -d '\n')
MARIADB_DATABASE = $(shell cat ../../../.env | grep MARIADB_DATABASE | awk -F'=' '{print $2}' | tr -d '\n')
MARIADB_USER = $(shell cat ../../../.env | grep MARIADB_USER | awk -F'=' '{print $2}' | tr -d '\n')
MARIADB_PASSWORD = $(shell cat ../../../.env | grep MARIADB_PASSWORD | awk -F'=' '{print $2}' | tr -d '\n')

run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 $(RABBITMQ_USER) $(RABBITMQ_PWD) S3_HOST=http://localhost:9000 $(S3_AUTH_KEY) $(S3_AUTH_PWD) $(USER_AUTH) $(PWD_AUTH) MARIADB_ADDR=localhost:3306 $(MARIADB_USER) $(MARIADB_PASSWORD) $(MARIADB_DATABASE) go run .
run-dev-remote:
	DEV_MODE=true USER_AUTH=$(USER_AUTH) PWD_AUTH=$(PWD_AUTH) go run .

build:
	go build -o build/api
build_image:
	docker build . -t voogle-api

generate-docs:
	swag init -g main.go
