HELM_VERSION = 3.10.2
HELM_ZIP = helm-v$(HELM_VERSION)-linux-amd64.tar.gz
HELM_URL = https://get.helm.sh/$(HELM_ZIP)

helm-install:		## Install Helm.
	@helm version | grep $(HELM_VERSION); \
		if [ $$? -eq 1 ]; then \
			wget -q $(HELM_URL); \
			tar zxvf $(HELM_ZIP); \
			sudo mv linux-amd64/helm /usr/local/bin/helm; \
			rm -rf $(HELM_ZIP) linux-amd64/ ; \
			helm version; \
		else \
			echo "${GREEN} [*] Helm already installed ${RESET}"; \
		fi

helm-add-repos:	## Add the required Helm Charts repositories.
	helm repo add aws-efs-csi-driver https://kubernetes-sigs.github.io/aws-efs-csi-driver/
	helm repo add couchbase https://couchbase-partners.github.io/helm-charts/
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo add jetstack https://charts.jetstack.io
	helm repo add metallb https://metallb.github.io/metallb
	helm repo add minio https://charts.min.io/
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add nsqio https://nsqio.github.io/helm-chart

	# Update your local Helm chart repository cache.
	helm repo update

helm-release:		## Install Helm release.
	make helm-add-repos
	make k8s-init-cert-manager
	make k8s-install-couchbase-crds
	cd $(ROOT_DIR)/deployments/saferwall \
		&& helm dependency update \
		&& helm install -name $(SAFERWALL_RELEASE_NAME) \
		 --namespace default .

helm-debug:		## Dry run install chart.
	cd $(ROOT_DIR)/deployments/saferwall \
		&& helm install -name $(SAFERWALL_RELEASE_NAME) \
			--debug --dry-run --namespace default . >> debug.yaml

helm-upgrade:		## Upgrade a given release.
	cd $(ROOT_DIR)/deployments/ \
		&& helm upgrade $(SAFERWALL_RELEASE_NAME) saferwall

helm-update-dep: # Update Helm deployement dependecies
	cd  $(ROOT_DIR)/deployments \
		&& helm dependency update saferwall
