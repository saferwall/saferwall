WINDOWS_DEFENDER_UPDATE = https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86
WINDOWS_DEFENDER_LOADLIBRARY_GIT = https://github.com/taviso/loadlibrary.git
WINDOWS_DEFENDER_INSTALL_DIR = /opt/windowsdefender

install-windefender:	## install Windows Defender
	sudo apt install gcc-multilib exiftool cabextract -y
	cd /tmp && git clone $(WINDOWS_DEFENDER_LOADLIBRARY_GIT)
	cd /tmp/loadlibrary && make
	sudo mv /tmp/loadlibrary $(WINDOWS_DEFENDER_INSTALL_DIR)
	sudo chown -R $(id -u):$(id -g) $(WINDOWS_DEFENDER_INSTALL_DIR)
	make update-windefender

update-windefender:		## update Windows Defender
	wget "$(WINDOWS_DEFENDER_UPDATE)" -O $(WINDOWS_DEFENDER_INSTALL_DIR)/engine/mpam-fe.exe
	cd $(WINDOWS_DEFENDER_INSTALL_DIR)/engine && cabextract mpam-fe.exe && rm mpam-fe.exe

uninstall-windefender:	## uninstall Windows Defender
	sudo rm -rf $(WINDOWS_DEFENDER_INSTALL_DIR)
