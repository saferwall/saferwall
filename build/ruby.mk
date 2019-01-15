install-ruby:
	sudo apt update
	curl -fsSL https://github.com/rbenv/rbenv-installer/raw/master/bin/rbenv-installer | bash
	echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.zshrc
	echo 'eval "$(rbenv init -)"' >> ~/.zshrc
	source ~/.zshrc
	curl -fsSL https://github.com/rbenv/rbenv-installer/raw/master/bin/rbenv-installer | bash
	rbenv install 2.4.1
	rbenv global 2.4.1
	ruby -v
	echo "gem: --no-document" > ~/.gemrc
	gem install bundler