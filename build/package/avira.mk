AVIRA_URL = http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
AVIRA_FUSEBUNDLE = http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
AVIRA_INSTALL_DIR = /opt/avira

install-avira:				## install Avira Linux Version
	wget $(AVIRA_URL) -P /tmp
	tar zxvf /tmp/scancl-linux_glibc22.tar.gz -C /tmp
	sudo mkdir -p /opt/avira
	sudo mv /tmp/scancl-1.9.161.2/* $(AVIRA_INSTALL_DIR)
	make update-avira

update-avira:				## update Avira Linux Version
	wget $(AVIRA_FUSEBUNDLE) -P /tmp
	unzip -o /tmp/avira_fusebundlegen-linux_glibc22-en.zip -d /tmp
	/tmp/fusebundle.bin
	sudo unzip -o /tmp/install/fusebundle-linux_glibc22-int.zip -d $(AVIRA_INSTALL_DIR)
	sudo rm -rf /tmp/*

uninstall-avira:			## uninstall Avira Linux Version
	sudo rm -rf $(AVIRA_INSTALL_DIR)