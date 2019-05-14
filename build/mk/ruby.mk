ruby-install:		## Install Ruby
	sudo apt-get update && sudo apt-get install ruby-full build-essential zlib1g-dev -y
	echo '# Install Ruby Gems to ~/gems' >> ~/.bashrc
	echo 'export GEM_HOME="$(HOME)/gems"' >> ~/.bashrc
	echo 'export PATH="$(HOME)/gems/bin:$(PATH)"' >> ~/.bashrc
	source ~/.bashrc

ruby-install-jekyll:	## Install Jekyll
	gem install jekyll bundler
