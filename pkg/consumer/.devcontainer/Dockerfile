# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.166.1/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then \
    su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

##### Install Prerequisites #####
RUN echo "Installing Prerequisites ..." \
    && apt-get update \
    && apt-get install -qq --no-install-recommends ca-certificates

##### Install dependencies #####
RUN echo "Installing Dependencies" \
    && buildDeps="automake unzip wget libtool make gcc pkg-config git" \
    && apt-get install -y $buildDeps --no-install-recommends

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
    && apt-get install -qq libimage-exiftool-perl

####### Installing File #######
RUN echo "Installing File..." \
    && apt-get install -qq file

####### Installing DiE #######
ENV DIE_VERSION     3.01
ENV DIE_URL         https://github.com/horsicq/DIE-engine/releases/download/$DIE_VERSION/die_lin64_portable_$DIE_VERSION.tar.gz
ENV DIE_ZIP         /tmp/die_lin64_portable_$DIE_VERSION.tar.gz
ENV DIE_DIR         /opt/die/

RUN echo "Installing DiE..." \
    && apt-get install libglib2.0-0 -y \
    && wget $DIE_URL -O $DIE_ZIP \
    && tar zxvf $DIE_ZIP -C /tmp \
    && mv /tmp/die_lin64_portable/ $DIE_DIR

####### Installing Yara #######
ENV YARA_VERSION        4.0.5
ENV YARA_ARCHIVE        $YARA_VERSION.tar.gz
ENV YARA_DOWNLOAD_URL   https://github.com/VirusTotal/yara/archive/v$YARA_ARCHIVE

RUN echo "Installing Yara..." \
    && apt-get install libssl-dev libglib2.0-0 -y \
    && wget $YARA_DOWNLOAD_URL \
    && tar zxvf v$YARA_ARCHIVE \
    && cd ./yara-$YARA_VERSION \
    && ./bootstrap.sh \
    && ./configure \
    && make \
    && make install \
    && rm -rf ./yara-$YARA_VERSION \
    && rm -f $YARA_ARCHIVE

####### Installing Capstone #######
RUN echo "Installing Capstone..." \
    && apt-get install libcapstone-dev -y

# Cleanup.
RUN rm -rf /tmp/* \
    && apt-get purge -y --auto-remove $buildDeps \
    && rm -rf /var/lib/apt/lists/*

