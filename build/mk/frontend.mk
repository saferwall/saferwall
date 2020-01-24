ui-docker-run:		## Run the docker container
	sudo docker run -it -p 80:80 --name ui saferwall/ui

frontend-build:		## Build the UI in docker
	sudo make docker-build IMG=ui VERSION=0.0.1 DOCKER_FILE=ui/Dockerfile DOCKER_DIR=ui/

frontend-release:		## build and release UI.
	sudo make docker-release IMG=ui VERSION=0.0.1 DOCKER_FILE=ui/Dockerfile DOCKER_DIR=ui/

	