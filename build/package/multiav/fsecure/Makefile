-include ../../common.mk

FSECURE_VERSION = 11.10.68
FSECURE_INSTALL_DIR = /opt/f-secure
FSECURE_UPDATE = http://download.f-secure.com/latest/fsdbupdate9.run
FSECURE_URL = https://download.f-secure.com/corpro/ls/trial/fsls-${FSECURE_VERSION}-rtm.tar.gz


install:	## install FSecure Linux Security
	apt install wget lib32stdc++6 rpm psmisc -y
	wget $(FSECURE_URL) -P /tmp
	tar zxvf /tmp/fsls-${FSECURE_VERSION}-rtm.tar.gz -C /tmp
	chmod a+x /tmp/fsls-${FSECURE_VERSION}-rtm/fsls-${FSECURE_VERSION}
	/tmp/fsls-${FSECURE_VERSION}-rtm/fsls-${FSECURE_VERSION} --auto standalone lang=en --command-line-only
	make update

update:		## update FSecure Linux Security
	wget $(FSECURE_UPDATE) -P /tmp
	mv /tmp/fsdbupdate9.run $(FSECURE_INSTALL_DIR)
	$(FSECURE_INSTALL_DIR)/fsav/bin/dbupdate $(FSECURE_INSTALL_DIR)/fsdbupdate9.run
	rm -rf /tmp/*

uninstall:	## uninstall FSecure Linux Security
	/opt/f-secure/fsav/bin/uninstall-fsav
	rm -rf /opt/f-secure/
