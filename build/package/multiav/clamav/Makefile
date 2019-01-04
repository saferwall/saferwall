-include ../../common.mk

install:	## install the Open Source ClamAV Antivirus
	apt-get install clamav clamav-daemon -y
	make update
	service clamav-daemon status
    
update:		## update the Open Source ClamAV Antivirus
	freshclam
	service clamav-daemon restart

uninstall:	## uninstall the Open Source ClamAV Antivirus
	apt-get remove clamav clamav-daemon -y
