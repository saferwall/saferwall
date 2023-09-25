
LUA_VER = 5.4.6
install-lua: ## Install lua
	sudo apt install -y build-essential libreadline-dev unzip wget
	curl -R -O http://www.lua.org/ftp/lua-${LUA_VER}.tar.gz
	tar -zxf lua-${LUA_VER}.tar.gz
	cd lua-${LUA_VER} \
		&& make linux test \
		&& sudo make install
	rm lua-${LUA_VER}.tar.gz && rm -rf lua-${LUA_VER}

LUAROCKS_VER = 3.9.1
install-luarocks: ## Install luarocks
	wget https://luarocks.org/releases/luarocks-${LUAROCKS_VER}.tar.gz
	tar zxpf luarocks-${LUAROCKS_VER}.tar.gz
	cd luarocks-${LUAROCKS_VER} \
		&& ./configure --lua-version=5.4 \
		&& make \
		&& sudo make install
	rm luarocks-${LUAROCKS_VER}.tar.gz
	rm -rf luarocks-${LUAROCKS_VER}

LUAROCKS_DIR=${ROOT_DIR}/services/sandbox/lua/.luarocks
install-lua-deps: ## Install lua dependencies
	luarocks install --tree ${LUAROCKS_DIR} luacov
	luarocks install --tree ${LUAROCKS_DIR} serpent
	luarocks install --tree ${LUAROCKS_DIR} luaunit
	luarocks install --tree ${LUAROCKS_DIR} argparse
	luarocks install --tree ${LUAROCKS_DIR} rapidjson
	luarocks install --tree ${LUAROCKS_DIR} luacheck
	luarocks install --tree ${LUAROCKS_DIR} luafilesystem
	luarocks install --tree ${LUAROCKS_DIR} dumbluaparser
	luarocks install --tree ${LUAROCKS_DIR} luacov-console
