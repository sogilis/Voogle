run:
	go run main.go

build:
	go build -o build/api

build_image:
	docker build . -t voogle-api

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

fmt:
	gofmt -s -w .