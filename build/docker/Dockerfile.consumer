FROM debian:buster-slim

LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="F-Secure for Linux in a docker container"


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

####### Installing Exiftool #######
RUN echo "Installing Exiftool..." \
    && apt-get install libimage-exiftool-perl -y

####### Installing Capstone #######
RUN echo "Installing Capstone..." \
    && apt-get install libcapstone-dev -y
