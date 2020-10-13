KIND_VERSION = v0.9.0

kind-install: ## Install Kind for local kubernetes cluster deployements.
	curl -Lo kind "https://kind.sigs.k8s.io/dl/$(KIND_VERSION)/kind-$$(uname)-amd64"
	chmod +x kind
	sudo cp kind /usr/local/bin && rm kind
	kind version

kind-create-cluster:	## Create Kind cluster.
	sudo kind get clusters
	sudo kind create cluster --name saferwall --config build/kind/cluster-config.yaml
	sudo kubectl cluster-info --context kind-saferwall

kind-deploy-ingress-nginx: ## Deploy ngress-nginx in Kind.
	sudo kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

kind-delete-cluster:	## Delete Kind cluster.
	sudo kind delete clusters saferwall

kind-up: kind-create-cluster kind-deploy-ingress-nginx helm-init-cert-manager		## Deploy Kind cluster and install requirements.
