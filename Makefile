lint-dockerfile:
	./tools/hadolint.sh

run-ci-locally:
	act

start_s3_local:
	docker run -d -v $$PWD/.data/minio:/home/shared -e MINIO_ROOT_USER=admin -e MINIO_ROOT_PASSWORD=password -p 9000:9000 -p 9001:9001 minio/minio:latest server --console-address ":9001" /home/shared
