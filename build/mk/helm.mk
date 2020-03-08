HELM_VERSION = 3.0.2
HELM_ZIP = helm-v$(HELM_VERSION)-linux-amd64.tar.gz 
HELM_URL = https://get.helm.sh/$(HELM_ZIP)

helm-install:		## Install Helm
	wget $(HELM_URL)
	tar zxvf $(HELM_ZIP)
	sudo mv linux-amd64/helm /usr/local/bin/helm
	rm -f $(HELM_ZIP)
	helm version

helm-init:			## Init Helm repo
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	# Make sure we get the latest list of charts
	helm repo update

helm-create:		## Create a helm deployment
	cd  $(ROOT_DIR)/deployments \
		&& helm create saferwall \ 
		&& helm ls

helm-release:		## Install Helm release
	cd  $(ROOT_DIR)/deployments \
		&& helm install saferwall --generate-name

helm-upgrade:		## Upgrade a given release
	helm upgrade $(RELEASE_NAME) saferwall
	