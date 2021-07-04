REPO = saferwall

docker-build: ## Build the container
	@docker build $(ARGS) -t $(REPO)/$(IMG) \
		-f $(DOCKER_FILE) $(DOCKER_DIR)

docker-build-nc: ## Build the container without caching
	@docker build ${ARGS} --no-cache -t $(REPO)/$(IMG) \
		-f $(DOCKER_FILE) $(DOCKER_DIR)

docker-run: ## Run container on port configured in `config.env`
	docker run -d -p 50051:50051 $(REPO)/$(IMG)

docker-up: build run ## Run container

docker-stop: ## Stop and remove a running container
	docker stop $(IMG); docker rm $(REPO)/$(IMG)

docker-release: docker-repo-login docker-build-nc docker-publish ## Make a release by building and publishing the `{version}` and `latest` tagged containers to ECR

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
	docker tag $(REPO)/$(IMG) $(REPO)/$(IMG):latest

docker-tag-version: 	## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(REPO)/$(IMG) $(REPO)/$(IMG):$(VERSION)

docker-repo-login: 	## Login to Docker Hub
	@echo 'login to docker hub'
	@echo '$(DOCKER_HUB_PWD)' | docker login -u $(DOCKER_HUB_USR) --password-stdin

docker-install:		## Install docker.
	sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo apt-key fingerprint 0EBFCD88
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $$(lsb_release -cs) stable"
	sudo apt-get update
	sudo apt-get install docker-ce -y

DOCKER_COMPOSE_VER=1.29.2
docker-compose-install: 	## Install docker-compose
	sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VER}/docker-compose-$$(uname -s)-$$(uname -m)" -o /usr/local/bin/docker-compose
	sudo chmod +x /usr/local/bin/docker-compose
	docker-compose --version

docker-stop-all:		## Stop all containers
	docker stop $$(docker ps -a -q)

docker-rm-all:			## Delete all containers
	docker rm $$(docker ps -a -q)

docker-rm-images:		## Delete all images
	docker rmi $$(docker images -q) -f

docker-rm-dangling:		## Delete all dangling images
	docker images --quiet --filter=dangling=true | xargs --no-run-if-empty docker rmi -f

docker-rm-image-tags:	## Delete all tags from image, {IMG} as argument.
	docker images | grep $(IMG) | tr -s ' ' | cut -d ' ' -f 2 | xargs -I {} docker rmi -f $(REPO)/$(IMG):{}

docker-get-ip:			## Get container IP addr
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(IMG)

docker-daemon-restart:	## Restart docker daemon & reload config
	sudo systemctl daemon-reload
	sudo systemctl restart docker

docker-enable-experimental:		## Enable experimental
	echo '{"experimental":true}' >> /etc/docker/daemon.json

docker-stats: ## Get docker stats nicely formatted.
	docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"

docker-non-root: ## Run docker as non root.
	sudo groupadd docker
	sudo usermod -aG docker $(USER)
	newgrp docker
	sudo chown "$(USER)":"$(USER)" /home/"$(USER)"/.docker -R
	sudo chmod g+rwx "$(HOME)/.docker" -R

docker-cadvisor:	## Run docker cAdvisor.
	docker run \
	--volume=/:/rootfs:ro \
	--volume=/var/run:/var/run:rw \
	--volume=/sys:/sys:ro \
	--volume=/var/lib/docker/:/var/lib/docker:ro \
	--volume=/dev/disk/:/dev/disk:ro \
	--publish=8080:8080 \
	--detach=true \
	--name=cadvisor \
	google/cadvisor:latest
