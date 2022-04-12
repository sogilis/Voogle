build:
	cargo build --release
run:
	RUST_LOG=trace cargo run

build-docker:
	docker build . -t env-generator

publish: build-docker
	docker tag env-generator ghcr.io/sogilis/env-generator:latest
	docker push ghcr.io/sogilis/env-generator:latest

