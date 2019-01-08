include ../../common.mk
include ../../vault.mk

AVIRA_URL = http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
AVIRA_FUSEBUNDLE = http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
AVIRA_INSTALL_DIR = /opt/avira

install:	## install Avira Linux Version
	make install-vault
	apt-get update
	apt-get install wget unzip libc6-i386 -y
	wget $(AVIRA_URL) -P /tmp
	tar zxvf /tmp/scancl-linux_glibc22.tar.gz -C /tmp
	mkdir -p /opt/avira
	mv /tmp/scancl-1.9.161.2/* $(AVIRA_INSTALL_DIR)
ifdef SAFERWALL_DEV
	cp hbedv.key $(AVIRA_INSTALL_DIR)
endif
ifdef SAFERWALL_TEST
	make install-vault
	vault read -field=avira multiav/avira | base64 -d > $(AVIRA_INSTALL_DIR)/hbedv.key
endif
	make update

update:		## update Avira Linux Version
	wget $(AVIRA_FUSEBUNDLE) -P /tmp
	unzip -o /tmp/avira_fusebundlegen-linux_glibc22-en.zip -d /tmp
	/tmp/fusebundle.bin
	unzip -o /tmp/install/fusebundle-linux_glibc22-int.zip -d $(AVIRA_INSTALL_DIR)
	rm -rf /tmp/*

uninstall:	## uninstall Avira Linux Version
	rm -rf $(AVIRA_INSTALL_DIR)
