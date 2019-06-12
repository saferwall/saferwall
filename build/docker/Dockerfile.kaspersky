FROM debian:stretch-slim
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Kaspersky Anti-Virus for Linux File Servers in a docker container"

# Vars

ENV KASPERSKY_VERSION       10.1.1.6421
ENV KASPERSKY_BIN 			/opt/kaspersky/kesl/bin/kesl-control
ENV KASPERSKY_SETUP 		/opt/kaspersky/kesl/bin/kesl-setup.pl
ENV KASPERSKY_URL 			https://products.s.kaspersky-labs.com/endpoints/keslinux10/10.1.1.6421/multilanguage-20190517_122450/02543683/kesl_10.1.1-6421_amd64.deb
ENV KASPERSKY_TMP			/tmp/kaspersky

# Install dependencies
RUN apt-get update \
	&& apt-get install wget perl locales procps -y \
	&& mkdir -p $KASPERSKY_TMP \
	&& wget -N $KASPERSKY_URL -P $KASPERSKY_TMP \
	&& dpkg -i $KASPERSKY_TMP/kesl_10.1.1-6421_amd64.deb

# Set up locales
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    locale-gen
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8   

# Configure it
ADD install.conf install.conf
RUN $KASPERSKY_SETUP --autoinstall=install.conf; exit 0

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN service kesl-supervisor start \
	&& $KASPERSKY_BIN --scan-file eicar \
	&& $KASPERSKY_BIN -E --query "EventType=='ThreatDetected'" \
	; exit 0

# Clean up
RUN rm -rf $KASPERSKY_TMP
