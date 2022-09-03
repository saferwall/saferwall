TRID_URL = http://mark0.net/download/trid_linux_64.zip
TRID_ZIP = /tmp/trid_linux_64.zip
TRID_DEFS_ZIP = /tmp/triddefs.zip
TRID_DEFS = http://mark0.net/download/triddefs.zip

trid-install:	## Install TRiD
	wget -N $(TRID_URL) -O $(TRID_ZIP)
	wget -N $(TRID_DEFS) -O $(TRID_DEFS_ZIP)
	unzip -o $(TRID_ZIP) -d /tmp
	unzip -o $(TRID_DEFS_ZIP) -d /tmp
	sudo mv /tmp/trid /usr/bin/
	sudo mv /tmp/triddefs.trd /usr/bin/
	chmod +x /usr/bin/trid
	## export LC_ALL=C
