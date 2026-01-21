couchabse/install:	## Install Couchbase in Ubuntu
	curl -O https://packages.couchbase.com/releases/couchbase-release/couchbase-release-1.0-amd64.deb
	sudo dpkg -i ./couchbase-release-1.0-amd64.deb
	sudo apt-get update
	sudo apt-get install couchbase-server -y

couchbase/start:	## Start Couchbase Server docker container.
	$(eval COUCHBASE_CONTAINER_STATUS := $(shell docker container inspect -f '{{.State.Status}}' $(COUCHBASE_CONTAINER_NAME)))
ifeq ($(COUCHBASE_CONTAINER_STATUS),running)
	@echo "All good, couchabse server is already running."
else ifeq ($(COUCHBASE_CONTAINER_STATUS),exited)
	@echo "Starting Couchbase Server ..."
	docker start $(COUCHBASE_CONTAINER_NAME)
else
	@echo "Creating Couchbase Server ..."
	docker run -d --name $(COUCHBASE_CONTAINER_NAME) -p 8091-8094:8091-8094 -p 11210:11210 $(COUCHBASE_CONTAINER_VER)
endif

couchbase/init:	## Init couchbase database by creating the cluster and required buckets.
	# Init the cluster.
	@echo "${GREEN} [*] =============== Creating Cluster =============== ${RESET}"
	docker compose start couchbase
	docker exec $(COUCHBASE_CONTAINER_NAME) \
		couchbase-cli cluster-init \
		--cluster localhost \
		--cluster-username $(COUCHBASE_ADMIN_USER) \
		--cluster-password $(COUCHBASE_ADMIN_PWD) \
		--cluster-port 8091 \
		--cluster-name saferwall \
		--services data,index,query \
		--cluster-ramsize 512 \
		--cluster-index-ramsize 256

	# Create require buckets.
	for bucket in $(COUCHBASE_BUCKETS_LIST) ; do \
		echo "${GREEN} [*] =============== Creating $$bucket =============== ${RESET}" ; \
		docker exec $(COUCHBASE_CONTAINER_NAME) \
			couchbase-cli bucket-create \
			--cluster localhost \
			--username $(COUCHBASE_ADMIN_USER) \
			--password $(COUCHBASE_ADMIN_PWD) \
			--bucket sfw \
			--bucket-type couchbase \
			--bucket-ramsize 128 \
			--max-ttl 500000000 \
			--enable-flush 0 ; \
	done
