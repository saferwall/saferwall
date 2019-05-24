FROM debian:stretch-slim AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Avira Linux Version in a docker container"

# Vars
ENV AVIRA_URL  http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
ENV AVIRA_FUSEBUNDLE http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
ENV AVIRA_INSTALL_DIR /opt/avira
ENV AVIRA_TMP /tmp/avira

# Install dependencies
RUN apt-get update \
    && apt-get install wget unzip libc6-i386 -y

# Install Avira
RUN wget $AVIRA_URL -P $AVIRA_TMP \
    && tar zxvf $AVIRA_TMP/scancl-linux_glibc22.tar.gz -C $AVIRA_TMP \
    && mkdir /opt/avira \
    && mv $AVIRA_TMP/scancl-1.9.161.2/* $AVIRA_INSTALL_DIR

# Update the VDF
RUN wget $AVIRA_FUSEBUNDLE -P $AVIRA_TMP \
	&& unzip -o $AVIRA_TMP/avira_fusebundlegen-linux_glibc22-en.zip -d $AVIRA_TMP \
	&& $AVIRA_TMP/fusebundle.bin \
	&& unzip -o $AVIRA_TMP/install/fusebundle-linux_glibc22-int.zip -d $AVIRA_INSTALL_DIR

# Apply the license
ADD hbedv.key $AVIRA_INSTALL_DIR

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN /opt/avira/scancl eicar ; exit 0

# Clean up
RUN rm -rf $AVIRA_TMP 
