HELM_VERSION = 2.14.3
HELM_ZIP = helm-v$(HELM_VERSION)-linux-amd64.tar.gz 
HELM_URL = https://get.helm.sh/$(HELM_ZIP)

helm-install:		## install vault
	wget $(HELM_URL)
	tar zxvf $(HELM_ZIP)
	sudo mv linux-amd64/helm /usr/local/bin/helm
	rm -f $(HELM_ZIP)