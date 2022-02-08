run:
	go run .

run-dev:
	DEV_MODE=true REDIS_ADDR=localhost:6379 REDIS_PWD= S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password go run .

run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/encoder

test:
	go test ./... -v

test-protobuf:
	cd contracts/v1; protoc --go_out=module=github.com/Sogilis/Voogle/services/encoder/contracts/v1:. *.proto

build_image:
	docker build . -t voogle-encoder

lint:
	golangci-lint run -v

lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

lint-protobuf_docker:
	docker run --rm -v $(PWD):/workspace -w /workspace bufbuild/buf lint

fmt:
	gofmt -s -w .
