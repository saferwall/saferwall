-include ../../common.mk

ESET_URL = https://download.eset.com/com/eset/apps/business/es/linux/latest/esets.amd64.deb.bin
ESET_LICENSE = ERA-Endpoint.lic
ESET_CONFIG_DIR = /etc/opt/eset
ESET_INSTALL_DIR = /opt/eset

install:	## install ESET File Server Security for Linux, Please Provide ESET_USER and ESET_PWD as arguments
	apt-get update
	apt-get install wget libc6-i386 -y
	wget $(ESET_URL) --user=$(ESET_USER) --password=$(ESET_PWD) -P /tmp 
	chmod +x /tmp/esets.amd64.deb.bin
	/tmp/esets.amd64.deb.bin --skip-license
	sed -i -e 's/#av_update_username = \"\"/av_update_username = \"$$(ESET_USER)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
	sed -i -e 's/#av_update_password = \"\"/av_update_password = \"$$(ESET_PWD)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
ifdef SAFERWALL_DEV
	cp $(ESET_LICENSE) $(ESET_CONFIG_DIR)/esets/license/ERA-Endpoint.lic
	$(ESET_INSTALL_DIR)/esets/sbin/esets_lic --import=$(ESET_INSTALL_DIR)/esets/etc/license/
endif
	rm -rf /tmp/*
	make update

update:		## update ESET File Server Security for Linux
	/opt/eset/esets/sbin/esets_update

uninstall:	## uninstall EST File Server Security for Linux
	dpkg --purge esets
