PACKER_VERSION	=	1.6.3
PACKER_ZIP 		= 	packer_$(PACKER_VERSION)_linux_amd64.zip
PACKER_URL 		= 	https://releases.hashicorp.com/packer/$(PACKER_VERSION)/$(PACKER_ZIP)


packer-install:		## Install packer from HashiCorp
	wget $(PACKER_URL)
	sudo unzip -o $(PACKER_ZIP) -d /usr/bin
	rm -f $(PACKER_ZIP)
	packer version
