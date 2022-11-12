local-debug:	## Locally start saferwall dependencies
	make couchbase-start
	make nsq-start
	make minio-start
