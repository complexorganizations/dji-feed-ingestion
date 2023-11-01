# Get the latest ubuntu image from dockerhub.
FROM ubuntu:latest

# Metadata as described above.
LABEL maintainer="support@parking-united.com" \
    version="1.0" \
    description="Docker image for DJI feed ingestion."

# Install the dependencies for the middleware; and then install the middleware.
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get dist-upgrade -y && \
    apt-get install curl -y && \
    curl https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/middleware/middleware-installer.sh -o middleware-installer.sh && \
    bash middleware-installer.sh && \
    apt-get autoremove -y && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/* middleware-installer.sh

# Open the ports for the RTMP. (DJI > RTMP > Server)
EXPOSE 1935/tcp 1935/udp
# Open the ports for the RTSP. (Server > RTSP > VLC)
EXPOSE 8554/tcp 8554/udp
