-include build/multiav/avast/Makefile
-include build/multiav/avira/Makefile
-include build/multiav/bitdefender/Makefile
-include build/multiav/clamav/Makefile
-include build/multiav/comodo/Makefile
-include build/multiav/eset/Makefile
-include build/multiav/fsecure/Makefile
-include build/multiav/kaspersky/Makefile
-include build/multiav/mcafee/Makefile

api:	## Generates protocol buffers definitions files. 
	protoc -I $(ROOT_DIR)/api/protobuf-spec/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:$(ROOT_DIR)/api/protobuf-spec/ \
		$(ROOT_DIR)/api/protobuf-spec/multiav.$(AV_VENDOR).proto

compile: api	## Compile gRPC server
	go build -ldflags "-s -w" -o bin/server server/main.go 

