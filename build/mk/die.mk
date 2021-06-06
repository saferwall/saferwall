DIE_VERSION = 2.05
DIE_URL = https://github.com/horsicq/DIE-engine/releases/download/$(DIE_VERSION)/die_lin64_portable_$(DIE_VERSION).tar.gz
DIE_ZIP = /tmp/die_lin64_portable_$(DIE_VERSION).tar.gz

die-install:	## Install DiE
	wget -N $(DIE_URL) -O $(DIE_ZIP)
	tar zxvf $(DIE_ZIP) -C /tmp
	sudo mv /tmp/die_lin64_portable/ /opt/die/
