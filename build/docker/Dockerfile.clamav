FROM alpine:3.9
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="The Open Source ClamAV Antivirus in a docker container"

# Install ClamAV
RUN apk --no-cache add clamav clamav-daemon

# Update ClamAV
RUN freshclam

# Add EICAR Anti-Virus Test File
ADD --chown=clamav:clamav http://www.eicar.org/download/eicar.com.txt eicar

# Configure permissions, this dir is required as this 
# is where clamd store its socket file
RUN mkdir /run/clamav  \
    && chown clamav:clamav /run/clamav

# Performs a simple test
RUN clamd && clamdscan eicar ; cat /var/log/clamav/clamd.log

# Cleanup
RUN rm -rf /var/lib/apt/lists/* && \
    rm eicar