jekyll-install:		## Install Jekyll
	sudo apt-get update && sudo apt-get install ruby-full build-essential zlib1g-dev -y
	echo '# Install Ruby Gems to ~/gems' >> ~/.bashrc
	echo 'export GEM_HOME="$(HOME)/gems"' >> ~/.bashrc
	echo 'export PATH="$(HOME)/gems/bin:$(PATH)"' >> ~/.bashrc
	source ~/.bashrc
	gem install jekyll bundler

jekyll-serve:		## Serve the jekyll app
	export LC_ALL=en_US.UTF-8
	export LANG=en_US.UTF-8
	bundle install
	bundle exec jekyll serve
