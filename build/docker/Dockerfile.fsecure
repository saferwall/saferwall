FROM debian:stretch-slim
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="FSecure Linux Security in a docker container"

# Vars
ENV FSECURE_VERSION     11.10.68
ENV FSECURE_INSTALL_DIR /opt/f-secure
ENV FSECURE_UPDATE      http://download.f-secure.com/latest/fsdbupdate9.run
ENV FSECURE_URL         https://download.f-secure.com/corpro/ls/trial/fsls-$FSECURE_VERSION-rtm.tar.gz
ENV FSECURE_TMP         /tmp/fsecure

# Install dependencies
RUN apt-get update \
    && apt-get install wget lib32stdc++6 rpm psmisc procps -y \
    && mkdir $FSECURE_TMP \
    && wget $FSECURE_URL -P $FSECURE_TMP \
	&& tar zxvf $FSECURE_TMP/fsls-$FSECURE_VERSION-rtm.tar.gz -C $FSECURE_TMP \
	&& chmod a+x $FSECURE_TMP/fsls-$FSECURE_VERSION-rtm/fsls-$FSECURE_VERSION \
	&& $FSECURE_TMP/fsls-$FSECURE_VERSION-rtm/fsls-$FSECURE_VERSION --auto standalone lang=en --command-line-only

# Update VDF
RUN wget $FSECURE_UPDATE -P $FSECURE_TMP \
	&& mv $FSECURE_TMP/fsdbupdate9.run $FSECURE_INSTALL_DIR \
	&& /etc/init.d/fsaua start \
	&& /etc/init.d/fsupdate start \
	&& $FSECURE_INSTALL_DIR/fsav/bin/dbupdate $FSECURE_INSTALL_DIR/fsdbupdate9.run ; exit 0 \
	&& /opt/f-secure/fsav/bin/fsav --version

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN /opt/f-secure/fsav/bin/fsav --virus-action1=report --suspected-action1=report eicar ; exit 0

# Clean up
RUN	rm -rf $FSECURE_TMP
