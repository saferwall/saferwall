ESET_URL 			= https://download.eset.com/com/eset/apps/business/es/linux/latest/esets.amd64.deb.bin
ESET_LICENSE 		= ERA-Endpoint.lic
ESET_CONFIG_DIR 	= /etc/opt/eset
ESET_INSTALL_DIR 	= /opt/eset
ESET_TEMP			= /tmp/eset

eset-install:		## install ESET File Server Security for Linux, ESET_USER / ESET_PWD are read from .env
	apt-get update
	apt-get install wget libc6-i386 ed -y
	mkdir -p $(ESET_TEMP)
	wget -N $(ESET_URL) --user=$(ESET_USER) --password=$(ESET_PWD) -P $(ESET_TEMP)
	chmod +x $(ESET_TEMP)/esets.amd64.deb.bin
	$(ESET_TEMP)/esets.amd64.deb.bin --skip-license
	sed -i -e 's/#av_update_username = \"\"/av_update_username = \"$(ESET_USER)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
	sed -i -e 's/#av_update_password = \"\"/av_update_password = \"$(ESET_PWD)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
ifdef SAFERWALL_DEV
	cp $(ROOT_DIR)/build/multiav/eset/$(ESET_LICENSE) $(ESET_CONFIG_DIR)/esets/license/ERA-Endpoint.lic
	$(ESET_INSTALL_DIR)/esets/sbin/esets_lic --import=$(ESET_INSTALL_DIR)/esets/etc/license/
endif
ifdef SAFERWALL_TEST
	vault read -field=ERA-Endpoint.lic multiav/eset | base64 -d > $(ESET_CONFIG_DIR)/esets/license/ERA-Endpoint.lic
	$(ESET_INSTALL_DIR)/esets/sbin/esets_lic --import=$(ESET_INSTALL_DIR)/esets/etc/license/
endif
	rm -rf $(ESET_TEMP)
	make update-eset

eset-update:		## update ESET File Server Security for Linux
	/opt/eset/esets/sbin/esets_update

eset-uninstall:		## uninstall EST File Server Security for Linux
	dpkg --purge esets
