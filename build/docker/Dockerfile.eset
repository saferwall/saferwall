FROM debian:stretch-slim
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="ESET File Server Security for Linux in a docker container"

# Vars
ARG var
ENV var=${var}


ENV ESET_URL 			https://download.eset.com/com/eset/apps/business/es/linux/latest/esets.amd64.deb.bin
ENV ESET_LICENSE 		ERA-Endpoint.lic
ENV ESET_CONFIG_DIR 	/etc/opt/eset
ENV ESET_INSTALL_DIR 	/opt/eset
ENV ESET_TEMP			/tmp/eset
ENV ESET_USER			user
ENV ESET_PWD			pass

# Install dependencies
RUN apt-get update \
    && apt-get install wget libc6-i386 ed -y

# Install ESET
RUN mkdir $ESET_TEMP \
	&& wget -N $ESET_URL --user=$ESET_USER --password=$ESET_PWD -P $ESET_TEMP \
	&& chmod +x $ESET_TEMP/esets.amd64.deb.bin \
	&& $ESET_TEMP/esets.amd64.deb.bin --skip-license

# Copy License Key
ADD ERA-Endpoint.lic $ESET_CONFIG_DIR/esets/license/ERA-Endpoint.lic

# Update the config
RUN	sed -i "s/#av_update_username = \"\"/av_update_username = \"$ESET_USER\"/g" $ESET_CONFIG_DIR/esets/esets.cfg \
	&& sed -i "s/#av_update_password = \"\"/av_update_password = \"$ESET_PWD\"/g" $ESET_CONFIG_DIR/esets/esets.cfg \
	&& $ESET_INSTALL_DIR/esets/sbin/esets_lic --import=$ESET_INSTALL_DIR/esets/etc/license/

# Update the VDF
RUN /opt/eset/esets/sbin/esets_update

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# Test detection
RUN /opt/eset/esets/sbin/esets_scan --clean-mode=NONE eicar ; exit 0

# Clean up
RUN rm -rf $ESET_TEMP
