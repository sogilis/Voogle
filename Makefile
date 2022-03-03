lint-dockerfile:
	./tools/hadolint.sh

run-ci-locally:
	act

start_all_services:
	docker-compose -f docker-compose-external.yml -f docker-compose-internal.yml up -d --build --remove-orphans;

start_external_services:
	docker-compose -f docker-compose-external.yml up -d --build --remove-orphans;

stop_services:
	docker-compose -f docker-compose-external.yml -f docker-compose-internal.yml stop;
	docker-compose -f docker-compose-external.yml -f docker-compose-internal.yml rm -f;