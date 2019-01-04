-include ../../common.mk

KASPERSKY_VERSION 		= 8.0.4-312
KASPERSKY_BIN 			= /opt/kaspersky/kav4fs/bin/kav4fs-control
KASPERSKY_SETUP 		= /opt/kaspersky/kav4fs/bin/kav4fs-setup.pl
KASPERSKY_LICENSE 		= /etc/kaspersky/license.key
KASPERSKY_INSTALL_DIR 	= /etc/kaspersky
KASPERSKY_URL 			= "https://products.s.kaspersky-labs.com/multilanguage/file_servers/kavlinuxserver8.0/kav4fs_${KASPERSKY_VERSION}_i386.deb"


install:	## install Kaspersky Anti-Virus for Linux File Servers
	apt-get update
	apt-get install wget libc6-i386 -y
	wget $(KASPERSKY_URL) -P /tmp
	dpkg -i /tmp/kav4fs_$(KASPERSKY_VERSION)_i386.deb
	chmod a+s $(KASPERSKY_BIN)
	mkdir -m 777 $(KASPERSKY_INSTALL_DIR)
	cp license.key $(KASPERSKY_LICENSE)
	$(KASPERSKY_BIN) --install-active-key $(KASPERSKY_LICENSE)
	$(KASPERSKY_SETUP) --auto-install install.conf
	rm -rf /tmp/*

update:		## update Kaspersky Anti-Virus for Linux File Servers
	$(KASPERSKY_BIN) --start-task 6
	$(KASPERSKY_BIN) --progress 6
	$(KASPERSKY_BIN) --get-stat Update

uninstall:	## uninstall Kaspersky Anti-Virus for Linux File Servers
	apt remove kav4fs -y
