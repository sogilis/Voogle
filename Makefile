lint-dockerfile:
	./tools/hadolint.sh

run-ci-locally:
	act

start_services:
	docker-compose up -d;

stop_services:
	docker-compose stop; docker-compose rm -f