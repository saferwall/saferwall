install-clamav:		## install the Open Source ClamAV Antivirus
	apt-get install clamav clamav-daemon -y
	make update-clamav
	service clamav-daemon status

update-clamav:		## update the Open Source ClamAV Antivirus
	freshclam
	service clamav-daemon restart

uninstall-clamav:	## uninstall the Open Source ClamAV Antivirus
	apt-get remove clamav clamav-daemon -y
