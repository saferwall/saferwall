make kind-install:
	curl -Lo kind "https://kind.sigs.k8s.io/dl/v0.8.1/kind-$$(uname)-amd64"
	chmod +x kind
	sudo cp kind /usr/local/bin && rm kind
	kind version
