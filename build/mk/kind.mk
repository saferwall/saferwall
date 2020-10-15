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

kind-deploy-ingress-nginx: ## Deploy ingress-nginx in Kind.
	# The manifests contains kind specific patches to forward the hostPorts to the ingress controller, 
	# set taint tolerations and schedule it to the custom labelled node.
	sudo kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
	# Wait a bit before probing for pods, otherwise you get: error: no matching resources found
	@echo "Sleeping 30 seconds" && sleep 30s
	# Now the Ingress is all setup. Wait until is ready to process requests running:
	sudo kubectl wait --namespace ingress-nginx \
	--for=condition=ready pod \
	--selector=app.kubernetes.io/component=controller \
	--timeout=90s

kind-delete-cluster:	## Delete Kind cluster.
	sudo kind delete clusters saferwall

kind-up: ## Deploy Kind cluster and install requirements.
	make kind-create-cluster
	make kind-deploy-ingress-nginx
