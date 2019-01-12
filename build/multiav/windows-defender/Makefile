WINDOWS_DEFENDER_UPDATE = https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86
WINDOWS_DEFENDER_LOADLIBRARY_GIT = https://github.com/taviso/loadlibrary.git
WINDOWS_DEFENDER_INSTALL_DIR = /opt/windowsdefender
WINDOWS_DEFENDER_TMP	= /tmp/windowsdefender

install-windowsdefender:	## install Windows Defender
	apt-get update
	apt-get install wget unzip libc6-i386 exiftool cabextract -y
	mkdir -p $(WINDOWS_DEFENDER_TMP)
	cd $(WINDOWS_DEFENDER_TMP) && git clone $(WINDOWS_DEFENDER_LOADLIBRARY_GIT)
	cd $(WINDOWS_DEFENDER_TMP)/loadlibrary && make
	mv $(WINDOWS_DEFENDER_TMP)/loadlibrary $(WINDOWS_DEFENDER_INSTALL_DIR)
	make update-windowsdefender

update-windowsdefender:		## update Windows Defender
	wget "$(WINDOWS_DEFENDER_UPDATE)" -O $(WINDOWS_DEFENDER_INSTALL_DIR)/engine/mpam-fe.exe
	cd $(WINDOWS_DEFENDER_INSTALL_DIR)/engine && cabextract mpam-fe.exe && rm mpam-fe.exe
	rm -rf $(WINDOWS_DEFENDER_TMP)

uninstall-windowsdefender:	## uninstall Windows Defender
	rm -rf $(WINDOWS_DEFENDER_INSTALL_DIR)
