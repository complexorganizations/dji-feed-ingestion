# Get the latest ubuntu image from docker hub
FROM ubuntu:latest

# Install the RTSP server with middleware.
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get dist-upgrade -y &&  \ 
    apt-get clean -y && \
    apt-get autoremove -y && \ 
    apt-get autoclean -y && \
    apt-get install -f -y && \
    apt-get install curl -y && \
    curl https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/middleware/middleware-installer.sh -o middleware-installer.sh && \
    bash middleware-installer.sh

# Open the ports for the RTMP. (DJI > RTMP > Server)
EXPOSE 1935/tcp
EXPOSE 1935/udp
# Open the ports for the RTSP. (Server > RTSP > VLC)
EXPOSE 8554/tcp
EXPOSE 8554/udp
