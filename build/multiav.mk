-include build/multiav/avast/Makefile

api:	## Generates protocol buffers definitions files. 
	protoc -I $(ROOT_DIR)/api/protobuf-spec/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:$(ROOT_DIR)/api/protobuf-spec/ \
		$(ROOT_DIR)/api/protobuf-spec/multiav.$(AV_VENDOR).proto

compile: api	## Compile gRPC server
	go build -ldflags "-s -w" -o bin/server server/main.go 

