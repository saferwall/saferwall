couchbase-run:		## Run couchbase docker container instance.
	sudo docker run -d --name db -p 8091-8094:8091-8094 -p 11210:11210 couchbase:enterprise-6.5.1

couchbase-start:	## Start exiting couchbase `db` container.
	sudo docker start db
