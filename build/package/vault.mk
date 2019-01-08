VAULT_ZIP = vault_1.0.1_linux_amd64.zip

install-vault:		## install vault
	wget https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_linux_amd64.zip
	sudo unzip -o $(VAULT_ZIP) -d /usr/bin
	rm -f $(VAULT_ZIP)