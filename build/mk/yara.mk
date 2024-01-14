YARA_VERSION = 4.4.0
YARA_ARCHIVE = ${YARA_VERSION}.tar.gz
YARA_DOWNLOAD_URL = https://github.com/VirusTotal/yara/archive/v${YARA_ARCHIVE}

yara-install:	# Install yara
	sudo apt-get install automake libtool make gcc pkg-config -y
	sudo apt-get install libssl-dev libglib2.0-0 libmagic-dev -y
	wget ${YARA_DOWNLOAD_URL}
	tar zxvf v${YARA_ARCHIVE}
	cd ./yara-${YARA_VERSION} \
		&& ./bootstrap.sh \
		&& ./configure --enable-magic \
		&& make -j $(shell nproc)\
		&& sudo make install \
		&& sudo ldconfig
	rm -rf ./yara-$(YARA_VERSION)
	rm -f $(YARA_ARCHIVE)
