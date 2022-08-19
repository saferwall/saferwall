NSQ_ZIP = nsq-1.2.1.linux-amd64.go1.16.6
NSQ_URL = https://s3.amazonaws.com/bitly-downloads/nsq/$(NSQ_ZIP).tar.gz

nsq-install:	## Install NSQ
	wget $(NSQ_URL) -P /tmp
	tar zxvf /tmp/$(NSQ_ZIP).tar.gz -C /tmp
	sudo mv /tmp/$(NSQ_ZIP)/bin/* /usr/bin
	rm -rf $(NSQ_ZIP)
	rm -f /tmp/$(NSQ_ZIP).tar.gz

nsq-start:		## Start nsqlookupd and nsqd locally
	nsqlookupd &
	nsqd --lookupd-tcp-address=127.0.0.1:4160 &
	nsqadmin --lookupd-http-address=127.0.0.1:4161 &

nsq-start-docker:	## Start nsq in docker
	docker pull nsqio/nsq
	docker run --name lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd
	ifconfig | grep addr
	# find ip then, For example, given a host IP of 172.17.42.1:
	docker run --name nsqd -p 4150:4150 -p 4151:4151 \
		nsqio/nsq /nsqd \
		--broadcast-address=172.17.42.1 \
		--lookupd-tcp-address=172.17.42.1:4160
