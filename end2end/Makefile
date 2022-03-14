default: test-dev

test-docker: start-docker test-dev clean-docker

start-docker:
	(cd .. && make start_all_services)

clean-docker:
	docker stop $$(docker ps -q)
	docker container rm $$(docker container ls -aq)
	docker volume rm -f $$(docker volume ls -q)

test-dev:
	E2E_API_ENDPOINT=http://localhost:4444 E2E_USER_NAME=dev E2E_USER_PWD=test go test ./... -v