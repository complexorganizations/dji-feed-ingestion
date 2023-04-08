# Get the latest ubuntu image from docker hub
FROM ubuntu:latest

# Install the RTSP server with middleware.
RUN apt-get update && \
    apt-get install curl -y && \
    curl https://raw.githubusercontent.com/complexorganizations/dji-feed-analysis/main/middleware/middleware-installer.sh -o middleware-installer.sh && \
    bash middleware-installer.sh

# Open the ports for the RTMP server. (DJI > RTMP > Server)
EXPOSE 1935/tcp
EXPOSE 1935/udp
# Open the ports for the RTSP server. (Server > RTSP > VLC)
EXPOSE 8554/tcp
EXPOSE 8554/udp
