MCAFEE_URL 			=	http://b2b-download.mcafee.com/products/evaluation/vcl/l64/vscl-l64-604-e.tar.gz
MCAFEE_UPDATE		= 	http://download.nai.com/products/DatFiles/4.x/nai/
MCAFEE_INSTALL_DIR 	= 	/opt/mcafee

install-mcafee:		## install McAfee VirusScan Command Line Scanner 
	wget $(MCAFEE_URL) -P /tmp
	sudo mkdir -p $(MCAFEE_INSTALL_DIR)
	sudo tar zxvf /tmp/vscl-l64-604-e.tar.gz -C $(MCAFEE_INSTALL_DIR)
	make update-mcafee

update-mcafee: 		## update McAfee VirusScan Command Line Scanner
	wget -Nc -r -nd -l1 -A "avvepo????dat.zip" http://download.nai.com/products/DatFiles/4.x/nai/ -P /tmp
	cd /tmp && unzip -o 'avvepo*'
	cd /tmp && sudo unzip -o 'avvdat-*' -d $(MCAFEE_INSTALL_DIR)
	$(MCAFEE_INSTALL_DIR)/uvscan --decompress
	sudo rm -rf /tmp/*

uninstall-mcafee:	## uninstall McAfee VirusScan Command Line Scanner
	echo 'y' | $(MCAFEE_INSTALL_DIR)/uninstall-uvscan
