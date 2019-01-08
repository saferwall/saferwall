# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help


# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


############################## DOCKER TASKS ##############################

build: ## Build the container
	docker build -t $(DOCKER_IMAGE) .

build-nc: ## Build the container without caching
	docker build --no-cache -t $(DOCKER_IMAGE) .

run: ## Run container on port configured in `config.env`
	docker run -d -p 50051:50051 $(DOCKER_IMAGE)

up: build run ## Run container

stop: ## Stop and remove a running container
	docker stop $(DOCKER_IMAGE); docker rm $(DOCKER_IMAGE)

release: build-nc publish ## Make a release by building and publishing the `{version}` ans `latest` tagged containers to ECR

publish: repo-login publish-latest publish-version ## Publish the `{version}` ans `latest` tagged containers to ECR

publish-latest: tag-latest ## Publish the `latest` taged container to ECR
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):latest

publish-version: tag-version ## Publish the `{version}` taged container to ECR
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

tag: tag-latest tag-version ## Generate container tags for the `{version}` ans `latest` tags

tag-latest: 	## Generate container `{version}` tag
	@echo 'create tag latest'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):latest

tag-version: 	## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

repo-login: 	## Login to Docker Hub
	docker login --username=$(DOCKER_HUB_USR) --password=$(DOCKER_HUB_PWD)

api:			## Generates protocol buffers definitions files. 
	protoc -I ../../../../api/protobuf-spec/ -I${GOPATH}/src --go_out=plugins=grpc:../../../../api/protobuf-spec/ ../../../../api/protobuf-spec/avast.proto

compile:
	go build -ldflags "-s -w" -o bin/server server/main.go 

VAULT_ZIP = vault_1.0.1_linux_amd64.zip
install-vault:		## install vault
	wget https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_linux_amd64.zip
	sudo unzip -o $(VAULT_ZIP) -d /usr/bin
	rm -f $(VAULT_ZIP)