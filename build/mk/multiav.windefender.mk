WINDOWS_DEFENDER_UPDATE = https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86
WINDOWS_DEFENDER_LOADLIBRARY = https://codeload.github.com/taviso/loadlibrary/zip/master
WINDOWS_DEFENDER_INSTALL_DIR = /opt/windowsdefender
WINDOWS_DEFENDER_TMP	= /tmp/windowsdefender

windefender-install:	## install Windows Defender
	apt-get update
	apt-get install wget unzip libc6-i386 gcc-multilib exiftool cabextract -y
	mkdir -p $(WINDOWS_DEFENDER_TMP)
	wget $(WINDOWS_DEFENDER_LOADLIBRARY) -P $(WINDOWS_DEFENDER_TMP)
	cd $(WINDOWS_DEFENDER_TMP) && unzip -o $(WINDOWS_DEFENDER_TMP)/master
	cd $(WINDOWS_DEFENDER_TMP)/loadlibrary-master && make
	mv $(WINDOWS_DEFENDER_TMP)/loadlibrary-master $(WINDOWS_DEFENDER_INSTALL_DIR)
	make windows-defender-update

windefender-update:		## update Windows Defender
	curl -sS -o $(WINDOWS_DEFENDER_INSTALL_DIR)/engine/mpam-fe.exe -L $(WINDOWS_DEFENDER_UPDATE)
	cd $(WINDOWS_DEFENDER_INSTALL_DIR)/engine && cabextract mpam-fe.exe && rm mpam-fe.exe
	rm -rf $(WINDOWS_DEFENDER_TMP)

windefender-uninstall:	## uninstall Windows Defender
	rm -rf $(WINDOWS_DEFENDER_INSTALL_DIR)
