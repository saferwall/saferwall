ui-docker-build:		## Build the UI in docker
	docker build -t saferwall/ui .


ui-docker-run:			## Run the docker container
	docker run -it -p 8080:80 --name ui saferwall/ui

ui: ui-docker-build ui-docker-run					## Build & run
	