run:
	go run .

run-dev:
	REDIS_ADDR=localhost:6379 REDIS_PWD= DEV_MODE=true S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password USER_AUTH=dev PWD_AUTH=test go run .

run-dev-remote:
	REDIS_ADDR=localhost:6379 REDIS_PWD= DEV_MODE=true USER_AUTH=dev PWD_AUTH=test go run .

build:
	go build -o build/api

test:
	go test ./... -v

test-protobuf:
	cd contracts/v1; protoc --go_out=module=github.com/Sogilis/Voogle/services/api/contracts/v1:. *.proto

build_image:
	docker build . -t voogle-api

lint:
	golangci-lint run -v

lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

lint-protobuf_docker:
	docker run --rm -v $(PWD):/workspace -w /workspace bufbuild/buf lint

fmt:
	gofmt -s -w .
