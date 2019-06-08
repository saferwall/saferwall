KASPERSKY_VERSION 		= 8.0.4-312
KASPERSKY_BIN 			= /opt/kaspersky/kav4fs/bin/kav4fs-control
KASPERSKY_SETUP 		= /opt/kaspersky/kav4fs/bin/kav4fs-setup.pl
KASPERSKY_LICENSE 		= /etc/kaspersky/license.key
KASPERSKY_INSTALL_DIR 	= /etc/kaspersky
KASPERSKY_URL 			= "https://products.s.kaspersky-labs.com/multilanguage/file_servers/kavlinuxserver8.0/kav4fs_${KASPERSKY_VERSION}_i386.deb"
KASPERSKY_TMP			= /tmp/kaspersky

kaspersky-install:	## install Kaspersky Anti-Virus for Linux File Servers
	sed -i 's/^# *\(en_US.UTF-8\)/\1/' /etc/locale.gen && locale-gen
	apt-get update
	apt-get install wget libc6-i386 -y
	mkdir -p $(KASPERSKY_TMP)
	wget -N $(KASPERSKY_URL) -P $(KASPERSKY_TMP)
	dpkg -i --force-architecture  $(KASPERSKY_TMP)/kav4fs_$(KASPERSKY_VERSION)_i386.deb
	chmod a+s $(KASPERSKY_BIN)
	mkdir -p -m 777 $(KASPERSKY_INSTALL_DIR)

ifdef SAFERWALL_DEV
	cp $(ROOT_DIR)/build/multiav/kaspersky/license.key $(KASPERSKY_LICENSE)
	$(KASPERSKY_BIN) --install-active-key $(KASPERSKY_LICENSE)
endif
ifdef SAFERWALL_TEST
	vault read -field=license.key multiav/kaspersky | base64 -d > $(KASPERSKY_LICENSE)
	$(KASPERSKY_BIN) --install-active-key $(KASPERSKY_LICENSE)
endif

	$(KASPERSKY_SETUP) --auto-install $(ROOT_DIR)/build/multiav/kaspersky/install.conf
	rm -rf $(KASPERSKY_TMP)

kaspersky-update:		## update Kaspersky Anti-Virus for Linux File Servers
	$(KASPERSKY_BIN) --start-task 6
	$(KASPERSKY_BIN) --progress 6
	$(KASPERSKY_BIN) --get-stat Update

kaspersky-uninstall:	## uninstall Kaspersky Anti-Virus for Linux File Servers
	apt remove kav4fs -y
