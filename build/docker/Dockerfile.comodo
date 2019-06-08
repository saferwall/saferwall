FROM debian:stretch-slim
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Comodo Antivirus for Linux in a docker container"

# Vars
ENV COMODO_URL          http://download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb
ENV COMODO_UPDATE       http://download.comodo.com/av/updates58/sigs/bases/bases.cav
ENV COMODO_INSTALL_DIR  /opt/COMODO


# Install dependencies
RUN apt-get update \
    && apt-get install wget binutils -y

# Install Comodo
RUN wget $COMODO_URL -P /tmp \
	&& cd /tmp && ar x cav-linux_x64.deb \
	&& tar zxvf /tmp/data.tar.gz -C /

# Update the VDF
ADD $COMODO_UPDATE /opt/COMODO/scanners/bases.cav

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN /opt/COMODO/cmdscan -v -s /eicar ; exit 0

# Clean up
RUN rm -f /tmp/*
