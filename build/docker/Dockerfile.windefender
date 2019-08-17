FROM debian:stretch-slim AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Windows Defender in a docker container"

# Vars
ENV WINDOWS_DEFENDER_UPDATE         https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86
ENV WINDOWS_DEFENDER_LOADLIBRARY    https://codeload.github.com/taviso/loadlibrary/zip/master
ENV WINDOWS_DEFENDER_INSTALL_DIR    /opt/windows-defender
ENV WINDOWS_DEFENDER_TMP            /tmp/windows-defender

# Install Windows Defender
RUN dpkg --add-architecture i386 \
    && apt-get update \ 
    && apt-get install wget unzip gcc-multilib exiftool cabextract build-essential libreadline-dev:i386 libc6-i386 -y \
    && mkdir $WINDOWS_DEFENDER_TMP \
    && wget $WINDOWS_DEFENDER_LOADLIBRARY -P $WINDOWS_DEFENDER_TMP \
    && cd $WINDOWS_DEFENDER_TMP \ 
    && unzip -o $WINDOWS_DEFENDER_TMP/master \
    && cd $WINDOWS_DEFENDER_TMP/loadlibrary-master \
    && make -j2 \
    && mv $WINDOWS_DEFENDER_TMP/loadlibrary-master $WINDOWS_DEFENDER_INSTALL_DIR

# Update the VDF
RUN apt-get install libreadline-dev -y && wget $WINDOWS_DEFENDER_UPDATE -O $WINDOWS_DEFENDER_INSTALL_DIR/engine/mpam-fe.exe \
    && cd $WINDOWS_DEFENDER_INSTALL_DIR/engine \
    && cabextract mpam-fe.exe \
    && rm mpam-fe.exe

# Add EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt eicar

# loadlibrary project from taviso is actually having issues, using local version instead
RUN rm -rf WINDOWS_DEFENDER_INSTALL_DIR
ADD windows-defender $WINDOWS_DEFENDER_INSTALL_DIR

# Performs a simple test
RUN cd $WINDOWS_DEFENDER_INSTALL_DIR && ./mpclient /eicar
