run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 RABBITMQ_USER=guest RABBITMQ_PWD=guest S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password USER_AUTH=dev PWD_AUTH=test MARIADB_ADDR=127.0.0.1:3306 MARIADB_USER=guest MARIADB_USER_PWD=guest MARIADB_NAME=voogle-database go run .
run-dev-remote:
	DEV_MODE=true USER_AUTH=dev PWD_AUTH=test go run .

build:
	go build -o build/api
build_image:
	docker build . -t voogle-api

generate-docs:
	swag init -g main.go
