test:
	go test ./... -v --cover

fmt:
	gofmt -s -w .

run-dev-api:
	(cd ./cmd/api && make run-dev)
run-dev-encoder:
	(cd ./cmd/encoder && make run-dev)
run-dev-gray-server-transformer:
	(cd ./cmd/gray-server-transformer && make run-dev)
run-dev-flip-server-transformer:
	(cd ./cmd/flip-server-transformer && make run-dev)

build-api:
	go build ./cmd/api
build-encoder:
	go build ./cmd/encoder
build-image-encoder:
	docker build -f cmd/encoder/Dockerfile . -t voogle-encoder
build-image-api:
	docker build -f cmd/api/Dockerfile . -t voogle-api

lint:
	golangci-lint run -v
lint_docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

generate-protobuf-contracts:
	cd pkg/contracts/v1; protoc --go_out=module=github.com/Sogilis/Voogle/src/pkg/contracts/v1:. *.proto

generate-protobuf-transformers:
	cd pkg/transformer; \
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative v1/transformer.proto;

lint-protobuf_docker:
	docker run --rm -v $(PWD):/workspace -w /workspace bufbuild/buf lint


