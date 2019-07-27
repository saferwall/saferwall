MINIO_URL = https://dl.min.io/server/minio/release/linux-amd64/minio

minio-install:	## Install minio locally
	wget $(MINIO_URL)
	sudo mv minio /usr/bin
	sudo chmod +x /usr/bin/minio

minio-start:  ## Start a minio server locally
	minio server --address localhost:9000 /tmp/minio

minio-docker-run:	## Run a mini docker container
	docker pull minio/minio
	docker run --name minio -it -p 9000:9000 minio/minio server /data

minio-docker-start:	## Start existing minio  container.
	docker start minio
