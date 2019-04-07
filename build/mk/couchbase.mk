couchbase-run:		## Run couchbase docker container instance.
	docker run -d --name db -p 8091-8094:8091-8094 -p 11210:11210 couchbase/server

couchbase-start:	## Run exiting couchbase `db` container.
	docker start db
