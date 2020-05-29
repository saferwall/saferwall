website-docker-run:		## Run the docker container
	sudo docker run -it -p 4000:4000 --name website saferwall/website

website-build: ## Build website in a docker container
	sudo make docker-build IMG=website DOCKER_FILE=build/docker/Dockerfile.website DOCKER_DIR=website/

website-release: ## Release website in a docker container
	sudo make docker-release IMG=website VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.website DOCKER_DIR=website/
