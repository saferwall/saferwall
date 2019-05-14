nodejs-install:		## Install nodejs
	curl -sL https://deb.nodesource.com/setup_8.x | sudo bash -
	sudo apt-get update
	sudo apt-get install -y nodejs
	node -v
	npm -v