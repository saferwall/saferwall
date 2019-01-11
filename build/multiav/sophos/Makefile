SOPHOS_INSTALL_ARCHIVE 	= 	sav-linux-free-9.tgz
SOPHOS_INSTALL_DIR 		= 	/opt/sophos
SOPHOS_INSTALL_ARGS 	=  --update-free --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s
SOPHOS_TMP				=   /tmp/sophos
SOPHOS_INSTALL_SCRIPT 	= 	$(SOPHOS_TMP)/sophos-av/install.sh

install-sophos:	## install Sophos Anti-Virus for Linux
	mkdir -p $(SOPHOS_TMP)
	tar zxvf $(SOPHOS_INSTALL_ARCHIVE) -C $(SOPHOS_TMP)
	$(SOPHOS_INSTALL_SCRIPT) $(SOPHOS_INSTALL_DIR) $(SOPHOS_INSTALL_ARGS)
	make update-sophos

update-sophos:		## update Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/update/savupdate.sh
	rm -rf $(SOPHOS_TMP)

uninstall-sophos:	## uninstall Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/uninstall.sh
