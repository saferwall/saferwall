-include ../../common.mk

SOPHOS_INSTALL_ARCHIVE 	= 	sav-linux-free-9.tgz
SOPHOS_INSTALL_SCRIPT 	= 	/tmp/sophos-av/install.sh
SOPHOS_INSTALL_DIR 		= 	/opt/sophos
SOPHOS_INSTALL_ARGS 	=  --update-free --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s

install:	## install Sophos Anti-Virus for Linux
	tar zxvf $(SOPHOS_INSTALL_ARCHIVE) -C /tmp/
	$(SOPHOS_INSTALL_SCRIPT) $(SOPHOS_INSTALL_DIR) $(SOPHOS_INSTALL_ARGS)
	make update

update:		## update Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/update/savupdate.sh
	rm -rf /tmp/*

uninstall:	## uninstall Anti-Virus for Linux
	$(SOPHOS_INSTALL_DIR)/uninstall.sh
