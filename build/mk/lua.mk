
LUA_VER = 5.4.6
install-lua: ## Install lua
	sudo apt install -y build-essential libreadline-dev unzip wget
	curl -R -O http://www.lua.org/ftp/lua-${LUA_VER}.tar.gz
	tar -zxf lua-${LUA_VER}.tar.gz
	cd lua-${LUA_VER} \
		&& make linux test \
		&& sudo make install
	rm lua-${LUA_VER}.tar.gz && rm -rf lua-${LUA_VER}
