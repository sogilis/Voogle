test:
	go test ./... -v

fmt:
	gofmt -s -w .

build-api:
	go build cmd/api
build-encoder:
	go build cmd/encoder

lint:
	golangci-lint run -v
lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

generate-protobuf:
	cd pkg/contracts/v1; protoc --go_out=module=github.com/Sogilis/Voogle/src/pkg/contracts/v1:. *.proto
lint-protobuf_docker:
	docker run --rm -v $(PWD):/workspace -w /workspace bufbuild/buf lint


