MCAFEE_URL 			=	http://b2b-download.mcafee.com/products/evaluation/vcl/l64/vscl-l64-604-e.tar.gz
MCAFEE_UPDATE		= 	http://download.nai.com/products/DatFiles/4.x/nai/
MCAFEE_INSTALL_DIR 	= 	/opt/mcafee
MCAFEE_TMP 			= 	/tmp/mcafee

install-mcafee:	## install McAfee VirusScan Command Line Scanner
	apt-get update
	apt-get install wget unzip -y
	mkdir -p $(MCAFEE_TMP)
	wget -N $(MCAFEE_URL) -P $(MCAFEE_TMP)
	mkdir -p $(MCAFEE_INSTALL_DIR)
	tar zxvf $(MCAFEE_TMP)/vscl-l64-604-e.tar.gz -C $(MCAFEE_INSTALL_DIR)
	make update-mcafee

update-mcafee:		## update McAfee VirusScan Command Line Scanner
	wget -Nc -r -nd -l1 -A "avvepo????dat.zip" http://download.nai.com/products/DatFiles/4.x/nai/ -P $(MCAFEE_TMP)
	cd $(MCAFEE_TMP) && unzip -o 'avvepo*'
	cd $(MCAFEE_TMP) && unzip -o 'avvdat-*' -d $(MCAFEE_INSTALL_DIR)
	$(MCAFEE_INSTALL_DIR)/uvscan --decompress
	rm -rf $(MCAFEE_TMP)

uninstall-mcafee:	## uninstall McAfee VirusScan Command Line Scanner
	echo 'y' | $(MCAFEE_INSTALL_DIR)/uninstall-uvscan
	rm -rf $(MCAFEE_INSTALL_DIR)
