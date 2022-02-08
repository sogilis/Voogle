run:
	go run .
run-dev:
	DEV_MODE=true REDIS_ADDR=localhost:6379 REDIS_PWD= S3_HOST=http://localhost:9000 S3_AUTH_KEY=admin S3_AUTH_PWD=password go run .
run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/encoder
build_image:
	docker build . -t voogle-encoder
