-include ../../common.mk

AVAST_GPG = http://files.avast.com/files/resellers/linux/avast.gpg
AVAST_LICENSE_FILE = license.avastlic


install:	## install Avast Security for Linux
	echo 'deb http://deb.avast.com/lin/repo debian release' | tee --append /etc/apt/sources.list
	apt-key adv --fetch-keys $(AVAST_GPG)
	apt-get update
	apt-get install avast -y
	touch /etc/avast/whitelist
	old='^DOWNLOAD=(.*)$$' && new='DOWNLOAD="curl -L -s -f"' \
		&& sed -i "s|$$old|$$new|g" /var/lib/avast/Setup/avast.setup
ifdef SAFERWALL_DEV
	@echo "Using Dev"
	cp $(AVAST_LICENSE_FILE) /etc/avast/license.avastlic
	chown avast:avast /etc/avast/license.avastlic
endif
ifdef SAFERWALL_TEST
	make install-vault
	vault read -field=avast multiav/avast | base64 -d > /etc/avast/license.avastlic
	chown avast:avast /etc/avast/license.avastlic
endif
	service avast start

update:		## update Avast Security for Linux
	/var/lib/avast/Setup/avast.vpsupdate

uninstall:	## uninstall Avast Security for Linux
	apt-get remove avast -y
