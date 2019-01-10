COMODO_URL = http://download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb
COMODO_UPDATE = http://download.comodo.com/av/updates58/sigs/bases/bases.cav
COMODO_INSTALL_DIR = /opt/COMODO
COMODO_LIB_SSL = http://security.ubuntu.com/ubuntu/pool/universe/o/openssl098/libssl0.9.8_0.9.8o-7ubuntu3.2.14.04.1_amd64.deb

install-comodo:	## install Comodo Antivirus for Linux
	apt-get update
	apt-get install wget -y
	wget $(COMODO_LIB_SSL) -P /tmp			# Download and install the trusty package manually
	wget $(COMODO_URL) -P /tmp
	dpkg -i /tmp/libssl0.9.8_0.9.8o-7ubuntu3.2.14.04.1_amd64.deb
	cd /tmp && ar x cav-linux_x64.deb
	tar zxvf /tmp/data.tar.gz -C /
	rm -f /tmp/cav-linux_x64.deb
	make update-comodo

update-comodo:		## update Comodo Antivirus for Linux
	wget -N $(COMODO_UPDATE) -P $(COMODO_INSTALL_DIR)/scanners/

uninstall-comodo:	## uninstall Comodo Antivirus for Linux
	rm -rf $(COMODO_INSTALL_DIR)
