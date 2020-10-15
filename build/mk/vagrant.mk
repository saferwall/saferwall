
VAGRANT_VERSION = 2.2.10
VAGRANT_ZIP_FILE = vagrant_$(VAGRANT_VERSION)_linux_amd64.zip
VAGRANT_DOWNLOAD_URL = https://releases.hashicorp.com/vagrant/$(VAGRANT_VERSION)/$(VAGRANT_ZIP_FILE)

vagrant-install: ## Download and install HashiCorp Vagrant.
	wget $(VAGRANT_DOWNLOAD_URL)
	sudo unzip -o $(VAGRANT_ZIP_FILE) -d /usr/bin
	rm -f $(VAGRANT_ZIP_FILE)
	vagrant version

vagrant-package: ## Package Vagrant box.
	$(eval VAGRANT_VM_NAME := $(shell VBoxManage list vms | cut -f 1 -d ' ' | tr -d '"'))
	vagrant package --base $(VAGRANT_VM_NAME) --output saferwall.box

vagrant-login:	## Authenticate to Vagrant cloud
	vagrant cloud auth login --token $(VAGRANT_TOKEN)

VAGRANT_DESCRIPTION = Saferwall kubernetes cluster for local use
VAGRANT_SHORT_DESCRIPTION = A hackable malware sandbox for the 21st Century
vagrant-publish:	## Upload the image to the cloud.
	vagrant cloud publish saferwall/saferwall $(SAFERWALL_VER) virtualbox saferwall.box \
	 -d "$(VAGRANT_DESCRIPTION)" --version-description "$(SAFERWALL_VER)"
	  --release --short-description "$(VAGRANT_SHORT_DESCRIPTION)" --force