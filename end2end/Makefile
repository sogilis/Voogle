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
	./launch.sh
