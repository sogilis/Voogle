run:
	go run .

run-dev:
	REDIS_ADDR=localhost:6666 REDIS_PWD= DEV_MODE=true S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password USER_AUTH=dev PWD_AUTH=test go run .

run-dev-remote:
	DEV_MODE=true USER_AUTH=dev PWD_AUTH=test go run .

build:
	go build -o build/api

test:
	go test ./... -v

build_image:
	docker build . -t voogle-api

lint:
	golangci-lint run -v

lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

fmt:
	gofmt -s -w .
