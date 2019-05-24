FROM debian:stretch-slim AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Bitdefender Scanner for Unices/Unix in a docker container"

# Vars
ENV BITDEFENDER_VERSION     7.7-1
ENV BITDEFENDER_ROOT_URL    http://download.bitdefender.com/SMB/Workstation_Security_and_Management/BitDefender_Antivirus_Scanner_for_Unices/Unix/Current/EN_FR_BR_RO/Linux
ENV BITDEFENDER_URL         $BITDEFENDER_ROOT_URL/BitDefender-Antivirus-Scanner-$BITDEFENDER_VERSION-linux-amd64.deb.run
ENV BITDEFENDER_INSTALLER   BitDefender-Antivirus-Scanner-$BITDEFENDER_VERSION-linux-amd64.deb.run
ENV BITDEFENDER_SCANNER     /opt/BitDefender-scanner/bin/bdscan

# Install dependencies
RUN apt-get update \
    && apt-get install wget psmisc -y

# Install Bitdefender
RUN wget $BITDEFENDER_URL -P /tmp \
	&& sed -i 's/^CRCsum=.*$/CRCsum="0000000000"/' /tmp/$BITDEFENDER_INSTALLER \
	&& sed -i 's/^MD5=.*$/MD5="00000000000000000000000000000000"/' /tmp/$BITDEFENDER_INSTALLER \
	&& sed -i 's/^more LICENSE$/cat  LICENSE/' /tmp/$BITDEFENDER_INSTALLER \
	&& chmod +x  /tmp/$BITDEFENDER_INSTALLER \
	&& (echo 'accept' ; echo 'n') | sh /tmp/$BITDEFENDER_INSTALLER --nox11

# Update the VDF
RUN bdscan --update

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN $BITDEFENDER_SCANNER eicar ; exit 0

# Clean up
RUN rm -rf /tmp/*
