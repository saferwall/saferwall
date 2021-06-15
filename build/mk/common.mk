protoc:
	protoc --go_out=plugins=grpc:. proto/avast.proto

build-d:
	go build -o server/server -ldflags "-s -w" server/main.go
	sudo docker build --no-cache -t $(docker_image) .

clean:
	sudo docker rmi -f $(docker images "dangling=true" -q)
	sudo docker rmi $(docker_image) -f

	docker rmi $(docker images --quiet --filter "dangling=true")


start:
	sudo docker run -d -p 50051:50051 $(docker_image)

compile:
	go build -o server/server -ldflags "-s -w" server/main.go



api:			## generates protocol buffers definitions files.
	protoc -I ../../../../api/protobuf-spec/ -I${GOPATH}/src --go_out=plugins=grpc:../../../../api/protobuf-spec/ ../../../../api/protobuf-spec/avast.proto


clean:
	# sudo docker rm $(docker ps -a -q)
	docker images --quiet --filter=dangling=true | xargs --no-run-if-empty docker rmi
	sudo docker rmi $(DOCKER_IMAGE) -f


compile:
	make api
	go build -ldflags "-s -w" -o bin/server server/main.go

buildxx: ## Build the container
	make api
	make compile
	sudo docker build -t $(DOCKER_IMAGE) .

build-ncxx: ## Build the container without caching
	make compile
	sudo docker build --no-cache -t $(DOCKER_IMAGE) .


upload:
	sudo apt install sharutils jq -y

	# AVAST
	cat license.avastlic | base64 | vault write "multiav/avast" license.avastlic=-
	vault read -field=license.avastlic multiav/avast | base64 -d > gpgkeybin2

	# AVIRA
	cat hbedv.key | base64 | vault write "multiav/avira" hbedv.key=-
	vault read -field=hbedv.key multiav/avira | base64 -d > gpgkeybin2

	# ESET
	vault kv put multiav/eset ESET_USER=EAV-0242798608 ESET_PWD=sv4rfnhvs6
	cat ERA-Endpoint.lic | base64 | vault write "multiav/eset" ERA-Endpoint.lic=-

	# KASPERSKY
	vault read -field=license.key multiav/kaspersky | base64 -d > $(KASPERSKY_LICENSE)
	cat license.key | base64 | vault write "multiav/kaspersky" license.key=-

	# Sophos
	cat sav-linux-free-9.tgz | base64 | vault write "multiav/sophos" sav-linux=-

	# ENV
	cat .env | base64 | vault write "secret/.env" .env=-
	vault read -field=.env secret/.env | base64 -d > .env

	# SAFERWALL BACKEND CONFIG FILE
	cat app.toml | base64 | vault write "secret/app.toml" app.toml=-