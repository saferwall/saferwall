FROM debian:stretch-slim
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Symantec Endpoint Protection Linux Client in a docker container"

# Vars
ENV SYMANTEC_DEB 	sep-deb.zip
ENV SYMANTEC_SAV	/opt/Symantec/symantec_antivirus/sav
ENV SYMANTEC_TMP    /tmp/symantec

# Install dependencies
RUN apt-get update \
    && apt-get install unzip kmod libc6-i386 -y

# Install Symantec
RUN mkdir -p $SYMANTEC_TMP
ADD $SYMANTEC_DEB $SYMANTEC_TMP
RUN unzip -o $SYMANTEC_TMP/$SYMANTEC_DEB -d $SYMANTEC_TMP \
    && $SYMANTEC_TMP/install.sh -i \
    && $SYMANTEC_SAV info --defs

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN /etc/init.d/symcfgd start \
    && /etc/init.d/rtvscand start \ 
    && /etc/init.d/smcd start \
    && /opt/Symantec/symantec_antivirus/sav manualscan --clscan /eicar \
    && TODAY=`date '+%m%d%Y'` ; cat /var/symantec/sep/Logs/$TODAY.log

# Clean up
RUN rm -rf $SYMANTEC_TMP