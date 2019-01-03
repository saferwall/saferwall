SYMANTEC_DEB 			= ./multiav/symantec/sep-deb.zip
SYMANTEC_SAV			= /opt/Symantec/symantec_antivirus/sav

install-symantec:	## install Symantec Endpoint Protection Linux Client
	unzip -o $(SYMANTEC_DEB) -d /tmp
	sudo sh /tmp/install.sh -i
	sudo rm -rf /tmp/*

update-symantec:	## update Symantec Endpoint Protection Linux Client
	$(SYMANTEC_SAV) liveupdate --update
	sudo $(SYMANTEC_SAV) --defs

uninstall-symantec:
	sudo /opt/Symantec/symantec_antivirus/uninstall.sh