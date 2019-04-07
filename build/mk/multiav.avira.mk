AVIRA_URL 			= http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
AVIRA_FUSEBUNDLE 	= http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
AVIRA_INSTALL_DIR 	= /opt/avira
AVIRA_TMP 			= /tmp/avira

install-avira:	## install Avira Linux Version
	apt-get update
	apt-get install wget unzip libc6-i386 -y
	wget $(AVIRA_URL) -P $(AVIRA_TMP)
	tar zxvf $(AVIRA_TMP)/scancl-linux_glibc22.tar.gz -C $(AVIRA_TMP)
	mkdir -p /opt/avira
	mv $(AVIRA_TMP)/scancl-1.9.161.2/* $(AVIRA_INSTALL_DIR)
ifdef SAFERWALL_DEV
	cp  $(ROOT_DIR)/build/multiav/avira/hbedv.key $(AVIRA_INSTALL_DIR)
endif
ifdef SAFERWALL_TEST
	vault read -field=hbedv.key multiav/avira | base64 -d > $(AVIRA_INSTALL_DIR)/hbedv.key
endif
	make update-avira

update-avira:		## update Avira Linux Version
	wget $(AVIRA_FUSEBUNDLE) -P $(AVIRA_TMP)
	unzip -o $(AVIRA_TMP)/avira_fusebundlegen-linux_glibc22-en.zip -d $(AVIRA_TMP)
	$(AVIRA_TMP)/fusebundle.bin
	unzip -o $(AVIRA_TMP)/install/fusebundle-linux_glibc22-int.zip -d $(AVIRA_INSTALL_DIR)
	rm -rf $(AVIRA_TMP)

uninstall-avira:	## uninstall Avira Linux Version
	rm -rf $(AVIRA_INSTALL_DIR)
