include ../../../.env

run:
	go run .
run-dev:
	DEV_MODE=true S3_HOST=http://localhost:9000 S3_AUTH_KEY=$(S3_AUTH_KEY) S3_AUTH_PWD=$(S3_AUTH_PWD) CONSUL_HOST=localhost:8500 go run .

run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/gray-server-transformer
build_image:
	docker build . -t voogle-gray-server-transformer
