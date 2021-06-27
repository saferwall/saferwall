AVAST_GPG = http://files.avast.com/files/resellers/linux/avast.gpg
AVAST_LICENSE_PATH = /etc/avast/license.avastlic

install-avast:	## install Avast Security for Linux
	echo 'deb http://deb.avast.com/lin/repo debian release' | tee --append /etc/apt/sources.list
	apt-key adv --fetch-keys $(AVAST_GPG)
	apt-get update
	apt-get install avast -y
	touch /etc/avast/whitelist
	old='^DOWNLOAD=(.*)$$' && new='DOWNLOAD="curl -L -s -f"' \
		&& sed -i "s|$$old|$$new|g" /var/lib/avast/Setup/avast.setup
ifdef SAFERWALL_DEV
	cp $(ROOT_DIR)/build/data/license.avastlic $(AVAST_LICENSE_PATH)
	chown avast:avast $(AVAST_LICENSE_PATH)
endif
ifdef SAFERWALL_TEST
	vault read -field=license.avastlic multiav/avast | base64 -d > $(AVAST_LICENSE_PATH)
	chown avast:avast $(AVAST_LICENSE_PATH)
endif
	service avast start

update-avast:		## update Avast Security for Linux
	/var/lib/avast/Setup/avast.vpsupdate

uninstall-avast:	## uninstall Avast Security for Linux
	apt-get remove avast -y