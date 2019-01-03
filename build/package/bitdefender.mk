BITDEFENDER_VERSION = 7.7-1
BITDEFENDER_ROOT_URL = http://download.bitdefender.com/SMB/Workstation_Security_and_Management/BitDefender_Antivirus_Scanner_for_Unices/Unix/Current/EN_FR_BR_RO/Linux
BITDEFENDER_URL = $(BITDEFENDER_ROOT_URL)/BitDefender-Antivirus-Scanner-$(BITDEFENDER_VERSION)-linux-amd64.deb.run
BITDEFENDER_INSTALLER = BitDefender-Antivirus-Scanner-${BITDEFENDER_VERSION}-linux-amd64.deb.run

install-bitdefender:		## install Bitdefender Scanner for Unices/Unix
	wget $(BITDEFENDER_URL) -P /tmp
	sed -i 's/^CRCsum=.*$$/CRCsum="0000000000"/' /tmp/$(BITDEFENDER_INSTALLER)
	sed -i 's/^MD5=.*$$/MD5="00000000000000000000000000000000"/' /tmp/$(BITDEFENDER_INSTALLER)
	sed -i 's/^more LICENSE$$/cat  LICENSE/' /tmp/$(BITDEFENDER_INSTALLER)
	chmod +x  /tmp/$(BITDEFENDER_INSTALLER)
	(echo 'accept' ; echo 'n') | sudo sh /tmp/$(BITDEFENDER_INSTALLER) --nox11
	make update-bitdefender

update-bitdefender:         ## update Bitdefender Scanner for Unices/Unix
	sudo bdscan --update
	sudo rm -rf /tmp/*

uninstall-bitdefender:		## uninstall Bitdefender Scanner for Unices/Unix
	echo 'y' | sudo sh /tmp/$(BITDEFENDER_INSTALLER) --uninstall --nox11
