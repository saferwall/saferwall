nodejs-install:		## Install nodejs
	curl -sL https://deb.nodesource.com/setup_8.x | sudo bash -
	apt-get update
	apt-get install -y nodejs
	node -v
	npm -v