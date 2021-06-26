KAFKA_VER = 2.8.0
KAFKA_PKG_NAME = kafka_2.13-$(KAFKA_VER)
KAFKA_INSTALL_DIR = /opt/kafka

kafka-setup:		## Downloads Kafka and save it.
	wget https://apache.mirror.digitalpacific.com.au/kafka/$(KAFKA_VER)/$(KAFKA_PKG_NAME).tgz
	tar -xzf $(KAFKA_PKG_NAME).tgz
	sudo mkdir -p $(KAFKA_INSTALL_DIR)
	sudo mv $(KAFKA_PKG_NAME)/* $(KAFKA_INSTALL_DIR)
	rm $(KAFKA_PKG_NAME).tgz && rm -r $(KAFKA_PKG_NAME)

kafka-start:		## Start Kafka broker service.
	# Start the ZooKeeper service and the Kafka broker service
	# Note: Soon, ZooKeeper will no longer be required by Apache Kafka.
	cd $(KAFKA_INSTALL_DIR) \
		&& sudo ./bin/zookeeper-server-start.sh config/zookeeper.properties &
		&& sleep 5s \
	 	&& sudo bin/kafka-server-start.sh config/server.properties