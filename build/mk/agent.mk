compile-agent:	## Build Windows Agent Server
		##GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o agent-server.exe cmd/agent/server/main.go
		mv agent-server.exe saferwall-agent-server-v$(SAFERWALL_VER)-windows-amd64.exe
		zip saferwall-agent-server-v$(SAFERWALL_VER)-windows-amd64.zip saferwall-agent-server-v$(SAFERWALL_VER)-windows-amd64.exe
