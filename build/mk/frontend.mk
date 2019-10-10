ui-docker-build:		## Build the UI in docker
	sudo docker build -t saferwall/ui .


ui-docker-run:			## Run the docker container
	sudo docker run -it -p 80:80 --name ui saferwall/ui

ui-docker-release:		## build and release UI.
	sudo make docker-release IMG=ui VERSION=0.0.1 DOCKER_FILE=ui/Dockerfile DOCKER_DIR=ui/

ui: ui-docker-build ui-docker-run					## Build & run
	