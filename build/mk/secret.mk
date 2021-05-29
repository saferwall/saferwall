gpg-encrypt-env:		## Encrypt .env with gpg.
	@gpg --cipher-algo AES256 --symmetric --batch --passphrase '$(SECRETS_PASSPHRASE)' .env
	mv .env.gpg $(ROOT_DIR)/.github/secrets

gpg-build-secrets:		## Encrypt secrets which contains license files with gpg.
	cd $(ROOT_DIR)/build/data \
		&& tar -P -cvzf secrets.tar.gz \
						ERA-Endpoint.lic \
						ERA-Endpoint.lic-expired \
						ESET_File_Security_for_Linux.lic \
						hbedv.key \
						hbedv.key.expired \
						kaspersky.license.key \
						kaspersky.license.key-expired \
						license.avastlic \
						license.avastlic.expired
	@gpg --cipher-algo AES256 --symmetric --batch \
		--passphrase '$(SECRETS_PASSPHRASE)' $(ROOT_DIR)/build/data/secrets.tar.gz
	mv $(ROOT_DIR)/build/data/secrets.tar.gz.gpg $(ROOT_DIR)/.github/secrets
