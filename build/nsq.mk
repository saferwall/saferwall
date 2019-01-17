NSQ_ZIP = nsq-1.1.0.linux-amd64.go1.10.3
NSQ_URL = https://s3.amazonaws.com/bitly-downloads/nsq/$(NSQ_ZIP).tar.gz

nsq-install:	## Install NSQ
	wget $(NSQ_URL) -P /tmp
	tar zxvf /tmp/$(NSQ_ZIP).tar.gz -C /tmp
	mv /tmp/$(NSQ_ZIP)/bin/* /usr/bin
	rm -rf $(NSQ_ZIP)
	rm -f /tmp/$(NSQ_ZIP).tar.gz