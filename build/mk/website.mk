website-build: ## Build website in a docker container
	sudo make docker-build IMG=website VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.website DOCKER_DIR=website/

website-release: ## Release website in a docker container
	sudo make docker-release IMG=website VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.website DOCKER_DIR=website/
