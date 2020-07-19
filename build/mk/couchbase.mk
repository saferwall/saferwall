COUCHBASE_CONTAINER_NAME = cb-db
COUCHBASE_CONTAINER_STATUS := $(shell sudo docker container inspect -f '{{.State.Status}}' $(COUCHBASE_CONTAINER_NAME))

couchbase-start:	## Start Couchbase Server docker container.
ifeq ($(COUCHBASE_CONTAINER_STATUS),running)
	@echo "All good, couchabse server is already running."
else
	@echo "No stress, creating one ..."
	sudo docker run -d --name $(COUCHBASE_CONTAINER_NAME) -p 8091-8094:8091-8094 -p 11210:11210 couchbase:enterprise-6.5.1
endif
