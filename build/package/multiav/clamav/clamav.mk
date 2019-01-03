install-clamav:         ## install the Open Source ClamAV Antivirus
	sudo apt install clamav clamav-daemon -y
	make update-clamav
	sudo service clamav-daemon status
    
update-clamav:          ## update the Open Source ClamAV Antivirus
	sudo freshclam
	sudo service clamav-daemon restart

uninstall-clamav:		## uninstall the Open Source ClamAV Antivirus
	sudo apt remove clamav clamav-daemon -y
