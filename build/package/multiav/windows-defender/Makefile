-include ../../common.mk

WINDOWS_DEFENDER_UPDATE = https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86
WINDOWS_DEFENDER_LOADLIBRARY_GIT = https://github.com/taviso/loadlibrary.git
WINDOWS_DEFENDER_INSTALL_DIR = /opt/windowsdefender

install:	## install Windows Defender
	apt-get update
	apt-get install wget unzip libc6-i386 exiftool cabextract -y
	cd /tmp && git clone $(WINDOWS_DEFENDER_LOADLIBRARY_GIT)
	cd /tmp/loadlibrary && make
	mv /tmp/loadlibrary $(WINDOWS_DEFENDER_INSTALL_DIR)
	make update

update:		## update Windows Defender
	wget "$(WINDOWS_DEFENDER_UPDATE)" -O $(WINDOWS_DEFENDER_INSTALL_DIR)/engine/mpam-fe.exe
	cd $(WINDOWS_DEFENDER_INSTALL_DIR)/engine && cabextract mpam-fe.exe && rm mpam-fe.exe

uninstall:	## uninstall Windows Defender
	rm -rf $(WINDOWS_DEFENDER_INSTALL_DIR)
