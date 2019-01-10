docker-build: ## Build the container
	docker build -t $(DOCKER_IMAGE) $(DOCKER_FILE)

docker-build-nc: ## Build the container without caching
	docker build --no-cache -t $(DOCKER_IMAGE) $(DOCKER_FILE)

docker-run: ## Run container on port configured in `config.env`
	docker run -d -p 50051:50051 $(DOCKER_IMAGE)

docker-up: build run ## Run container

docker-stop: ## Stop and remove a running container
	docker stop $(DOCKER_IMAGE); docker rm $(DOCKER_IMAGE)

docker-release: docker-build-nc docker-publish ## Make a release by building and publishing the `{version}` ans `latest` tagged containers to ECR

docker-publish: docker-repo-login docker-publish-latest docker-publish-version ## Publish the `{version}` ans `latest` tagged containers to ECR

docker-publish-latest: docker-tag-latest ## Publish the `latest` taged container to ECR
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):latest

docker-publish-version: docker-tag-version ## Publish the `{version}` taged container to ECR
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

docker-tag: docker-tag-latest docker-tag-version ## Generate container tags for the `{version}` ans `latest` tags

docker-tag-latest: 	## Generate container `{version}` tag
	@echo 'create tag latest'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):latest

docker-tag-version: 	## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

docker-repo-login: 	## Login to Docker Hub
	docker login --username=$(DOCKER_HUB_USR) --password=$(DOCKER_HUB_PWD)
