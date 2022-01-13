run:
	go run .

run-dev:
	USER_AUTH=dev PWD_AUTH=test go run .

build:
	go build -o build/api

test:
	go test ./...

build_image:
	docker build . -t voogle-api

lint:
	golangci-lint run -v

lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

fmt:
	gofmt -s -w .