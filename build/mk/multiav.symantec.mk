SYMANTEC_DEB 	= sep-deb.zip
SYMANTEC_SAV	= /opt/Symantec/symantec_antivirus/sav
SYMANTEC_TMP	= /tmp/symantec

symantec-install:	## install Symantec Endpoint Protection Linux Client
	apt-get update
	apt-get install unzip libc6-i386 -y
	mkdir -p $(SYMANTEC_TMP)
ifdef SAFERWALL_DEV
	cp $(ROOT_DIR)/build/multiav/symantec/$(SYMANTEC_DEB) $(SYMANTEC_TMP)
endif
ifdef SAFERWALL_TEST
	wget $(SYMANTEC_URL) -P $(SYMANTEC_TMP)
endif
	unzip -o $(SYMANTEC_TMP)/$(SYMANTEC_DEB) -d $(SYMANTEC_TMP)
	$(SYMANTEC_TMP)/install.sh -i
	echo "Sleeping 4 minutes for updates to be applied"
	sleep 4m
	$(SYMANTEC_SAV) info --defs
	rm -rf $(SYMANTEC_TMP)

symantec-update:		## update Symantec Endpoint Protection Linux Client
	$(SYMANTEC_SAV) liveupdate --update
	$(SYMANTEC_SAV) info --defs

symantec-uninstall:	## uninstall Symantec Endpoint Protection Linux Client
	echo 'Y' | /opt/Symantec/symantec_antivirus/uninstall.sh
