USER_AUTH = $(shell cat ../../../.env | grep USER_AUTH | awk -F'=' '{print $2}' | tr -d '\n')
PWD_AUTH = $(shell cat ../../../.env | grep PWD_AUTH | awk -F'=' '{print $2}' | tr -d '\n')

default: test-dev

test-docker: start-docker test-dev clean-docker

start-docker:
	(cd .. && make start_all_services)
	docker ps -a

clean-docker:
	docker stop $$(docker ps -q)
	docker container rm $$(docker container ls -aq)
	docker volume rm -f $$(docker volume ls -q)
	
test-dev:
	E2E_API_ENDPOINT=http://localhost:4444 E2E_USER_NAME=$(USER_AUTH) E2E_USER_PWD=$(PWD_AUTH) go test ./... -v