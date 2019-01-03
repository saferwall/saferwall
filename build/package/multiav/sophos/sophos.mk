SOPHOS_INSTALL_ARCHIVE 	= 	./multiav/sophos/sav-linux-free-9.tgz
SOPHOS_INSTALL_SCRIPT 	= 	/tmp/sophos-av/install.sh
SOPHOS_INSTALL_DIR 		= 	/opt/sophos
SOPHOS_INSTALL_ARGS 	=  --update-free --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s

install-sophos:			## install Sophos Anti-Virus for Linux
	tar zxvf $(SOPHOS_INSTALL_ARCHIVE) -C /tmp/
	sudo $(SOPHOS_INSTALL_SCRIPT) $(SOPHOS_INSTALL_DIR) $(SOPHOS_INSTALL_ARGS)
	make update-sophos

update-sophos:			## update Anti-Virus for Linux
	sudo $(SOPHOS_INSTALL_DIR)/update/savupdate.sh
	sudo rm -rf /tmp/*

uninstall-sophos:			## uninstall Anti-Virus for Linux
	sudo $(SOPHOS_INSTALL_DIR)/uninstall.sh
