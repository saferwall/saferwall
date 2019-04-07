SOPHOS_INSTALL_DIR 		= 	/opt/sophos
SOPHOS_INSTALL_ARGS 	=  --update-free --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s
SOPHOS_TMP				=   /tmp/sophos
SOPHOS_INSTALL_SCRIPT 	= 	$(SOPHOS_TMP)/sophos-av/install.sh
SOPHOS_INSTALL_ARCHIVE 	= 	$(SOPHOS_TMP)/sav-linux-free

install-sophos:	## install Sophos Anti-Virus for Linux
	mkdir -p $(SOPHOS_TMP)
ifdef SAFERWALL_DEV
	cp  $(ROOT_DIR)/build/multiav/sophos/sav-linux-free-9.tgz $(SOPHOS_INSTALL_ARCHIVE)
endif
ifdef SAFERWALL_TEST
	wget $(SOPHOS_URL) -P $(SOPHOS_TMP)
endif
	tar zxvf $(SOPHOS_INSTALL_ARCHIVE) -C $(SOPHOS_TMP)
	$(SOPHOS_INSTALL_SCRIPT) $(SOPHOS_INSTALL_DIR) $(SOPHOS_INSTALL_ARGS)
	make update-sophos

update-sophos:		## update Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/update/savupdate.sh
	rm -rf $(SOPHOS_TMP)

uninstall-sophos:	## uninstall Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/uninstall.sh
