lint-dockerfile:
	./tools/hadolint.sh

run-ci-locally:
	act

start_services: start_s3 start_redis
stop_services: stop_s3 stop_redis

start_s3:
	docker run -d --name voogle-s3 -v $$PWD/.data/minio:/home/shared -e MINIO_ROOT_USER=admin -e MINIO_ROOT_PASSWORD=password -p 9000:9000 -p 9001:9001 minio/minio:latest server --console-address ":9001" /home/shared

stop_s3:
	docker stop `docker ps -aqf "name=voogle-s3"`; docker rm voogle-s3

start_redis:
	docker run --name redis-test-instance -p 6379:6379 -d redis

stop_redis:
	docker stop redis-test-instance; docker rm redis-test-instance
