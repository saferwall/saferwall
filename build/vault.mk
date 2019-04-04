VAULT_VERSION = 1.0.3
VAULT_ZIP = vault_$(VAULT_VERSION)_linux_amd64.zip
VAULT_URL = https://releases.hashicorp.com/vault/$(VAULT_VERSION)/$(VAULT_ZIP)

install-vault:		## install vault
	wget $(VAULT_URL)
	sudo unzip -o $(VAULT_ZIP) -d /usr/bin
	rm -f $(VAULT_ZIP)