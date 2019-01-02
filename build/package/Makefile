# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

AVAST_GPG = http://files.avast.com/files/resellers/linux/avast.gpg
AVAST_LICENSE_FILE = ./multiav/avast/license.avastlic

install-avast:				## install Avast Security for Linux
	echo 'deb http://deb.avast.com/lin/repo debian release' | sudo tee --append /etc/apt/sources.list
	sudo apt-key adv --fetch-keys $(AVAST_GPG)
	sudo apt update
	sudo apt install avast -y
	sudo cp $(AVAST_LICENSE_FILE) /etc/avast/license.avastlic
	sudo chown avast:avast /etc/avast/license.avastlic
	/etc/init.d/avast status

update-avast:  				## update Avast Security for Linux
	sudo /var/lib/avast/Setup/avast.vpsupdate

uninstall-avast:			## uninstall Avast Security for Linux
	sudo apt remove avast -y



AVIRA_URL = http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
AVIRA_FUSEBUNDLE = http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
AVIRA_INSTALL_DIR = /opt/avira

install-avira:				## install Avira Linux Version
	wget $(AVIRA_URL) -P /tmp
	tar zxvf /tmp/scancl-linux_glibc22.tar.gz -C /tmp
	sudo mkdir -p /opt/avira
	sudo mv /tmp/scancl-1.9.161.2/* $(AVIRA_INSTALL_DIR)
	make update-avira

update-avira:				## update Avira Linux Version
	wget $(AVIRA_FUSEBUNDLE) -P /tmp
	unzip -o /tmp/avira_fusebundlegen-linux_glibc22-en.zip -d /tmp
	/tmp/fusebundle.bin
	sudo unzip -o /tmp/install/fusebundle-linux_glibc22-int.zip -d $(AVIRA_INSTALL_DIR)
	sudo rm -rf /tmp/*

uninstall-avira:			## uninstall Avira Linux Version
	sudo rm -rf $(AVIRA_INSTALL_DIR)
