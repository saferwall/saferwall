build-consumer:		## build consumer docker container
	sudo make docker-release IMG=consumer VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=.
