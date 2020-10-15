
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
	vagrant package --base $(VAGRANT_VM_NAME) --output saferwall

