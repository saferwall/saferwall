dist = $$(lsb_release -cs)

vbox-install:		## Install VirtualBox
	@vboxmanage --version | grep 6.1; \
		if [ $$? -eq 1 ]; then \
			wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add - ; \
			wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add - ; \
			sudo add-apt-repository "deb [arch=amd64] http://download.virtualbox.org/virtualbox/debian $(dist) contrib" ; \
			sudo apt-get update ; \
			sudo apt-get install virtualbox-6.1 -y ; \
			sudo apt-get install -f ; \
		else \
            echo "${GREEN} [*] VirtualBox already installed ${RESET}"; \
		fi

vbox-troubleshoot:
	vboxmanage startvm <vm-uuid> --type emergencystop

