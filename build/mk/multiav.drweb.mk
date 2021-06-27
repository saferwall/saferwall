DR_WEB_MAJOR_VERSION = 11.1
DR_WEB_MINOR_VERSION = 1
DR_WEB_INSTALLER = drweb-$(DR_WEB_MAJOR_VERSION).$(DR_WEB_MINOR_VERSION)-av-linux-amd64.run
DR_WEB_URL = https://download.geo.drweb.com/pub/drweb/unix/workstation/$(DR_WEB_MAJOR_VERSION)/$(DR_WEB_INSTALLER)

install-drweb:	## install DrWeb Antivirus for Linux
	sudo apt-get update
	sudo apt-get install wget netbase -y
	wget $(DR_WEB_URL) -P /tmp
	chmod +x /tmp/$(DR_WEB_INSTALLER)
	sudo /tmp/$(DR_WEB_INSTALLER) -- --non-interactive
	sudo /opt/drweb.com/bin/drweb-configd -d
	/opt/drweb.com/bin/drweb-ctl license --GetRegistered $(DR_WEB_LICENSE_KEY)
	make update-drweb
	sudo rm -rf /tmp/*

update-drweb:		## update DrWeb Antivirus for Linux
	sudo /opt/drweb.com/bin/drweb-configd -d  \
		&& /opt/drweb.com/bin/drweb-ctl update &> /dev/null; exit 0 \
		&& @echo "Updating the database ..." \
		&& /bin/bash -c 'while /opt/drweb.com/bin/drweb-ctl baseinfo | grep -q "Last successful update: unknown"; do sleep 5; done' \
		&& /opt/drweb.com/bin/drweb-ctl baseinfo \
		&& /opt/drweb.com/bin/drweb-ctl appinfo

uninstall-drweb:	## uninstall DrWeb Antivirus for Linux
	yes | sudo /opt/drweb.com/bin/uninst.sh
