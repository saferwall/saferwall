KIND_VERSION = 0.17.0
kind-install: ## Install Kind for local kubernetes cluster deployements.
	curl -o kind -sS -L https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-linux-amd64
	chmod +x kind
	sudo mv kind /usr/local/bin
	kind version

KIND_CLUSTER_NAME = sfw
kind-create-cluster:	## Create Kind cluster.
	kind get clusters
	kind create cluster --name $(KIND_CLUSTER_NAME) --config build/k8s/kind-cluster-config.yaml
	kubectl cluster-info --context kind-sfw

kind-deploy-ingress-nginx: ## Deploy ingress-nginx in Kind.
	# The manifests contains kind specific patches to forward the hostPorts to the ingress controller,
	# set taint tolerations and schedule it to the custom labelled node.
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
	# Wait a bit before probing for pods, otherwise you get: error: no matching resources found
	sleep 30s
	# Now the Ingress is all setup. Wait until is ready to process requests running:
	kubectl wait --namespace ingress-nginx \
	--for=condition=ready pod \
	--selector=app.kubernetes.io/component=controller \
	--timeout=90s
	kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission

kind-down:	## Delete Kind cluster.
	kind delete clusters ${KIND_CLUSTER_NAME}

kind-up: ## Deploy Kind cluster and install requirements.
	make kind-down || true
	make kind-create-cluster
	make kind-deploy-ingress-nginx
