-include ../../common.mk

SYMANTEC_DEB 			= ./multiav/symantec/sep-deb.zip
SYMANTEC_SAV			= /opt/Symantec/symantec_antivirus/sav

install:	## install Symantec Endpoint Protection Linux Client	
	apt-get update
	apt-get install unzip -y
	unzip -o $(SYMANTEC_DEB) -d /tmp
	sh /tmp/install.sh -i
	rm -rf /tmp/*

update:		## update Symantec Endpoint Protection Linux Client
	$(SYMANTEC_SAV) liveupdate --update
	$(SYMANTEC_SAV) --defs

uninstall:	## uninstall Symantec Endpoint Protection Linux Client
	/opt/Symantec/symantec_antivirus/uninstall.sh
