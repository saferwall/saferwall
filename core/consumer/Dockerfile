FROM ubuntu:bionic

##### Install Prerequisites #####
RUN echo "Updating packages ..." \ 
    && apt-get update \
    && apt-get install wget unzip -y

######## Installing TRiD ########
RUN echo "Installing TRiD..." \
    && wget http://mark0.net/download/trid_linux_64.zip -O /tmp/trid_linux_64.zip \
    && wget http://mark0.net/download/triddefs.zip -O /tmp/triddefs.zip \
    && cd /tmp \
    && unzip trid_linux_64.zip \
    && unzip triddefs.zip \
    && chmod +x trid \
    && mv trid /usr/bin/ \
    && mv triddefs.trd /usr/bin/
    # trid issue a sigfault if export LC_ALL=C not set

####### Installing Exiftool #######
RUN echo "Installing Exiftool..." \
    && apt-get install libimage-exiftool-perl -y

####### Installing Capstone #######
RUN echo "Installing Capstone..." \
    && apt-get install libcapstone-dev -y

####### Installing File #######
RUN echo "Installing File..." \
    && apt-get install file -y

####### Installing DiE #######
ENV DIE_VERSION     2.03
ENV DIE_URL         https://github.com/horsicq/DIE-engine/releases/download/$DIE_VERSION/die_lin64_portable_$DIE_VERSION.tar.gz
ENV DIE_ZIP         /tmp/die_lin64_portable_$DIE_VERSION.tar.gz
ENV DIE_DIR         /opt/die/

RUN echo "Installing DiE..." \
    && apt-get install libglib2.0-0 -y \
    && wget $DIE_URL -O $DIE_ZIP \
	&& tar zxvf $DIE_ZIP -C /tmp \
	&& mv /tmp/die_lin64_portable/ $DIE_DIR

WORKDIR /app

# Copy our static executable.
COPY consumer consumer
COPY saferwall.toml saferwall.toml

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN chown -R saferwall:saferwall .
RUN chmod +x consumer
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $DIE_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["./consumer"]