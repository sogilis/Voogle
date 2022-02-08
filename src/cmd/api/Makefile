run:
	go run .
run-dev:
	REDIS_ADDR=localhost:6379 REDIS_PWD= DEV_MODE=true S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password USER_AUTH=dev PWD_AUTH=test go run .
run-dev-remote:
	REDIS_ADDR=localhost:6379 REDIS_PWD= DEV_MODE=true USER_AUTH=dev PWD_AUTH=test go run .

build:
	go build -o build/api
build_image:
	docker build . -t voogle-api
