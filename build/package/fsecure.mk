
FSECURE_VERSION = 11.10.68
FSECURE_INSTALL_DIR = /opt/f-secure
FSECURE_UPDATE = http://download.f-secure.com/latest/fsdbupdate9.run
FSECURE_URL = https://download.f-secure.com/corpro/ls/trial/fsls-${FSECURE_VERSION}-rtm.tar.gz


install-fsecure:	## install FSecure Linux Security
	sudo apt install wget lib32stdc++6 rpm psmisc -y
	wget $(FSECURE_URL) -P /tmp
	tar zxvf /tmp/fsls-${FSECURE_VERSION}-rtm.tar.gz -C /tmp
	chmod a+x /tmp/fsls-${FSECURE_VERSION}-rtm/fsls-${FSECURE_VERSION}
	sudo /tmp/fsls-${FSECURE_VERSION}-rtm/fsls-${FSECURE_VERSION} --auto standalone lang=en --command-line-only
	make update-fsecure

update-fsecure:		## update FSecure Linux Security
	wget $(FSECURE_UPDATE) -P /tmp
	sudo mv /tmp/fsdbupdate9.run $(FSECURE_INSTALL_DIR)
	sudo $(FSECURE_INSTALL_DIR)/fsav/bin/dbupdate $(FSECURE_INSTALL_DIR)/fsdbupdate9.run
	sudo rm -rf /tmp/*

uninstall-fsecure:	## uninstall FSecure Linux Security
	sudo /opt/f-secure/fsav/bin/uninstall-fsav
	sudo rm -rf /opt/f-secure/
