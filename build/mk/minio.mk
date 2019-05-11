MINIO_URL = https://dl.min.io/server/minio/release/linux-amd64/minio

minio-install:	## Install minio locally
	wget $(MINIO_URL)
	mv minio /usr/bin
	chmod +x /usr/bin/minio

minio-docker:	## Install minio in docker
	docker pull minio/minio
	docker run -it -p 9000:9000 minio/minio server /data
