YARA_VERSION = 4.3.2
YARA_ARCHIVE = ${YARA_VERSION}.tar.gz
YARA_DOWNLOAD_URL = https://github.com/VirusTotal/yara/archive/v${YARA_ARCHIVE}
YARA_REPO_REPO  = https://github.com/saferwall/yara-rules.git
YARA_RULES_DIR  = /opt/yara-rules

# export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
# echo "/usr/local/lib" >> /etc/ld.so.conf
# ldconfig
yara-install:	# Install yara
	sudo apt-get install automake libtool make gcc pkg-config -y
	sudo apt-get install libssl-dev libglib2.0-0 -y
	wget ${YARA_DOWNLOAD_URL}
	tar zxvf v${YARA_ARCHIVE}
	cd ./yara-${YARA_VERSION} \
		&& ./bootstrap.sh \
		&& ./configure \
		&& make \
		&& sudo make install \
		&& sudo ldconfig
	rm -rf ./yara-$(YARA_VERSION)
	rm -f $(YARA_ARCHIVE)
	sudo git clone $(YARA_REPO_REPO) $(YARA_RULES_DIR)
