KASPERSKY_VERSION 		= 8.0.4-312
KASPERSKY_BIN 			= /opt/kaspersky/kav4fs/bin/kav4fs-control
KASPERSKY_SETUP 		= /opt/kaspersky/kav4fs/bin/kav4fs-setup.pl
KASPERSKY_LICENSE 		= /etc/kaspersky/license.key
KASPERSKY_INSTALL_CONG 	= ./multiav/kaspersky/install.conf
KASPERSKY_INSTALL_DIR 	= /etc/kaspersky
KASPERSKY_URL 			= "https://products.s.kaspersky-labs.com/multilanguage/file_servers/kavlinuxserver8.0/kav4fs_${KASPERSKY_VERSION}_i386.deb"


install-kaspersky:		## install Kaspersky Anti-Virus for Linux File Servers
	sudo apt install gcc-multilib
	wget $(KASPERSKY_URL) -P /tmp
	sudo dpkg -i /tmp/kav4fs_$(KASPERSKY_VERSION)_i386.deb
	sudo chmod a+s $(KASPERSKY_BIN)
	sudo mkdir -m 777 $(KASPERSKY_INSTALL_DIR)
	sudo cp ./multiav/kaspersky/license.key $(KASPERSKY_LICENSE)
	$(KASPERSKY_BIN) --install-active-key $(KASPERSKY_LICENSE)
	sudo $(KASPERSKY_SETUP) --auto-install $(KASPERSKY_INSTALL_CONG)
	sudo rm -rf /tmp/*

update-kaspersky:		## update Kaspersky Anti-Virus for Linux File Servers
	$(KASPERSKY_BIN) --start-task 6
	$(KASPERSKY_BIN) --progress 6
	$(KASPERSKY_BIN) --get-stat Update

uninstall_kaspersky:	## uninstall Kaspersky Anti-Virus for Linux File Servers
	sudo apt remove kav4fs -y
