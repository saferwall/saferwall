ESET_URL = https://download.eset.com/com/eset/apps/business/es/linux/latest/esets.amd64.deb.bin
ESET_LICENSE = ./multiav/eset/ERA-Endpoint.lic
ESET_CONFIG_DIR = /etc/opt/eset
ESET_INSTALL_DIR = /opt/eset

install-eset:		## install ESET File Server Security for Linux, Please Provide ESET_USER and ESET_PWD as arguments
	sudo apt install gcc-multilib -y
	wget $(ESET_URL) --user=$(ESET_USER) --password=$(ESET_PWD) -P /tmp 
	chmod +x /tmp/esets.amd64.deb.bin
	sudo /tmp/esets.amd64.deb.bin --skip-license
	sudo sed -i -e 's/#av_update_username = \"\"/av_update_username = \"$$(ESET_USER)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
	sudo sed -i -e 's/#av_update_password = \"\"/av_update_password = \"$$(ESET_PWD)\"/' $(ESET_CONFIG_DIR)/esets/esets.cfg
	sudo cp $(ESET_LICENSE) $(ESET_CONFIG_DIR)/esets/license/ERA-Endpoint.lic
	sudo $(ESET_INSTALL_DIR)/esets/sbin/esets_lic --import=$(ESET_INSTALL_DIR)/esets/etc/license/
	make update-eset

update-etset:		## update ESET File Server Security for Linux
	sudo /opt/eset/esets/sbin/esets_update

uninstall-eset:		## uninstall EST File Server Security for Linux
	sudo dpkg --purge esets
