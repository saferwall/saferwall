dist = $$(lsb_release -cs)

vbox-install:		## Install VirtualBox
	wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
	wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -
	sudo add-apt-repository "deb [arch=amd64] http://download.virtualbox.org/virtualbox/debian $(dist) contrib"
	sudo apt-get update
	sudo apt-get install virtualbox-6.1 -y
	sudo apt-get install -f
