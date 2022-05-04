lint-dockerfile:
	./tools/hadolint.sh

run-ci-locally:
	act

start_all_services:
	docker-compose --env-file .env -f docker/docker-compose-external.yml -f docker/docker-compose-internal.yml up -d --build --remove-orphans;

start_all_services_and_observability:
	docker-compose --env-file .env -f docker/docker-compose-external.yml -f docker/docker-compose-internal.yml -f docker/docker-compose-observability.yml up -d --build --remove-orphans;

start_external_services:
	docker-compose --env-file .env -f docker/docker-compose-external.yml up -d --build --remove-orphans;

stop_services:
	docker-compose -f docker/docker-compose-external.yml -f docker/docker-compose-internal.yml -f docker/docker-compose-observability.yml stop;
	docker-compose -f docker/docker-compose-external.yml -f docker/docker-compose-internal.yml -f docker/docker-compose-observability.yml rm -f;

E2E-tests: generate-env-file
	(cd end2end && make test-docker)

generate-env-file:
	docker run -v $$PWD:/env ghcr.io/sogilis/env-generator
