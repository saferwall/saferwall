dist = $$(lsb_release -cs)

vbox-install:		## Install VirtualBox
	@vboxmanage --version ; \
		if [ $$? -eq 1 ]; then \
			wget -O- -q https://www.virtualbox.org/download/oracle_vbox_2016.asc | sudo gpg --dearmour -o /usr/share/keyrings/oracle_vbox_2016.gpg ; \
			echo "deb [arch=amd64 signed-by=/usr/share/keyrings/oracle_vbox_2016.gpg] http://download.virtualbox.org/virtualbox/debian $(dist) contrib" | sudo tee /etc/apt/sources.list.d/virtualbox.list ; \
			sudo apt-get update ; \
			sudo apt-get install virtualbox-7.0 -y ; \
			sudo apt-get install -f ; \
		else \
			echo "${GREEN} [*] VirtualBox already installed ${RESET}"; \
		fi

vbox-troubleshoot:
	vboxmanage startvm <vm-uuid> --type emergencystop

vbox-install-fedora:
	sudo dnf -y install @development-tools
	sudo dnf -y install kernel-headers kernel-devel dkms elfutils-libelf-devel qt5-qtx11extras
	wget http://download.virtualbox.org/virtualbox/rpm/fedora/virtualbox.repo
	sudo mv virtualbox.repo /etc/yum.repos.d/virtualbox.repo
	sudo dnf install VirtualBox-6.1 -y
	sudo usermod -a -G vboxusers $USER
