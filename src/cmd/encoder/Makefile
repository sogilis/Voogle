RABBITMQ_USER = $(shell cat ../../../.env | grep RABBITMQ_USER | awk -F'=' '{print $2}' | tr -d '\n')
RABBITMQ_PWD = $(shell cat ../../../.env | grep RABBITMQ_PWD | awk -F'=' '{print $2}' | tr -d '\n')
S3_AUTH_KEY = $(shell cat ../../../.env | grep S3_AUTH_KEY | awk -F'=' '{print $2}' | tr -d '\n')
S3_AUTH_PWD = $(shell cat ../../../.env | grep S3_AUTH_PWD | awk -F'=' '{print $2}' | tr -d '\n')

run:
	go run .
run-dev:
	DEV_MODE=true RABBITMQ_ADDR=localhost:5672 $(RABBITMQ_USER) $(RABBITMQ_PWD) S3_HOST=http://localhost:9000 $(S3_AUTH_KEY) $(S3_AUTH_PWD) go run .
run-dev-remote:
	DEV_MODE=true go run .

build:
	go build -o build/encoder
build_image:
	docker build . -t voogle-encoder
