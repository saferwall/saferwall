MINIO_URL = https://dl.min.io/server/minio/release/linux-amd64/minio

minio-install:	## Install minio locally
	wget $(MINIO_URL)
	sudo mv minio /usr/bin
	sudo chmod +x /usr/bin/minio

minio-start:  ## Start a minio server locally
	sudo mkdir -p /samples
	sudo chmod -R 777 /samples
	minio server --address localhost:9000 ./minio

minio-docker-run:	## Run a mini docker container
	sudo docker pull minio/minio
	sudo docker run --name minio  \
		-e "MINIO_ACCESS_KEY=minio" \
  		-e "MINIO_SECRET_KEY=minio123" \
  		-it -p 9000:9000 minio/minio server /data

minio-docker-start:	## Start existing minio  container.
	sudo docker start minio
