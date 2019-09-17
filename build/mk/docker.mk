REPO = saferwall

docker-build: ## Build the container
	docker build $(ARGS) -t $(REPO)/$(IMG) -f $(DOCKER_FILE) $(DOCKER_DIR)

docker-build-nc: ## Build the container without caching
	docker build ${ARGS} --no-cache -t $(REPO)/$(IMG) -f $(DOCKER_FILE) $(DOCKER_DIR)

docker-run: ## Run container on port configured in `config.env`
	docker run -d -p 50051:50051 $(REPO)/$(IMG)

docker-up: build run ## Run container

docker-stop: ## Stop and remove a running container
	docker stop $(IMG); docker rm $(REPO)/$(IMG)

docker-release: docker-build docker-publish ## Make a release by building and publishing the `{version}` and `latest` tagged containers to ECR

docker-publish: docker-repo-login docker-publish-latest docker-publish-version ## Publish the `{version}` and `latest` tagged containers to ECR

docker-publish-latest: docker-tag-latest ## Publish the `latest` taged container to ECR
	@echo 'publish latest to $(REPO)/$(IMG)'
	docker push $(REPO)/$(IMG):latest

docker-publish-version: docker-tag-version ## Publish the `{version}` taged container to ECR
	@echo 'publish $(VERSION) to $(IMG)'
	docker push $(REPO)/$(IMG):$(VERSION)

docker-tag: docker-tag-latest docker-tag-version ## Generate container tags for the `{version}` and `latest` tags

docker-tag-latest: 	## Generate container `{version}` tag
	@echo 'create tag latest'
	docker tag $(REPO)/$(IMG) $(IMG):latest

docker-tag-version: 	## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(REPO)/$(IMG) $(REPO)/$(IMG):$(VERSION)

docker-repo-login: 	## Login to Docker Hub
	docker login --username=$(REPO) --password=$(DOCKER_HUB_PWD)

docker-install:		## install docker
	sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo apt-key fingerprint 0EBFCD88
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $$(lsb_release -cs) stable"
	sudo apt-get update
	sudo apt-get install docker-ce -y

docker-stop-all:		## Stop all containers
	docker stop $$(docker ps -a -q)

docker-rm-all:			## Delete all containers
	sudo docker rm $$(sudo docker ps -a -q)

docker-rm-images:		## Delete all images
	docker rmi $$(docker images -q) -f

docker-rm-dangling:		## Delete all dangling images
	docker images --quiet --filter=dangling=true | xargs --no-run-if-empty docker rmi -f

docker-rm-image-tags:	## Delete all tags from image
	docker images | grep $(IMG) | tr -s ' ' | cut -d ' ' -f 2 | xargs -I {} docker rmi -f $(REPO)/$(IMG):{}

docker-get-ip:			## Get container IP addr
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(IMG)

docker-daemon-restart:	## Restart docker daemon & reload config
	sudo systemctl daemon-reload
	sudo systemctl restart docker

docker-enable-experimental:		## Enable experimental
	echo '{"experimental":true}' >> /etc/docker/daemon.json
