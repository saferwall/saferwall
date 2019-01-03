install-clamav:         ## install the Open Source ClamAV Antivirus
	sudo apt install clamav clamav-daemon -y
	make update-clamav
    
update-clamav:          ## update the Open Source ClamAV Antivirus
	sudo freshclam

uninstall-clamav:		## uninstall the Open Source ClamAV Antivirus
	sudo apt remove clamav clamav-daemon -y
