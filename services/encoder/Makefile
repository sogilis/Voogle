run:
	go run .

run-dev:
	DEV_MODE=true go run .

run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/api

test:
	go test ./... -v

build_image:
	docker build . -t voogle-encoder

lint:
	golangci-lint run -v

lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

fmt:
	gofmt -s -w .
