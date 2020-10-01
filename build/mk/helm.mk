HELM_VERSION = 3.3.4
HELM_ZIP = helm-v$(HELM_VERSION)-linux-amd64.tar.gz 
HELM_URL = https://get.helm.sh/$(HELM_ZIP)

helm-install:		## Install Helm.
	wget $(HELM_URL)
	tar zxvf $(HELM_ZIP)
	sudo mv linux-amd64/helm /usr/local/bin/helm
	rm -f $(HELM_ZIP)
	helm version

helm-add-repos:	## Add the required Helm Charts repositories.
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm repo add couchbase https://couchbase-partners.github.io/helm-charts/
	helm repo add elastic https://helm.elastic.co
	helm repo add jetstack https://charts.jetstack.io
	# Update your local Helm chart repository cache.
	helm repo update

helm-create:		## Create a Helm deployment.
	cd $(ROOT_DIR)/deployments \
		&& helm create saferwall \ 
		&& helm ls

helm-release:		## Install Helm release.
	cd $(ROOT_DIR)/deployments \
		&& helm install saferwall --generate-name

helm-upgrade:		## Upgrade a given release.
	helm upgrade $(RELEASE_NAME) saferwall

helm-init-cert-manager: # Init cert-manager
	# Create the namespace for cert-manager.
	kubectl create namespace cert-manager
	# Install the CustomResourceDefinition
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.2/cert-manager.crds.yaml

helm-update-dependency: # Update Helm deployement dependecies
	cd  $(ROOT_DIR)/deployments \
		&& helm dependency update saferwall