# Create a RTSP server.

```
apt-get update
apt-get upgrade -y
apt-get install ffmpeg wget coreutils screen ffmpeg -y
mkdir /etc/rtsp-simple-server
cd /etc/rtsp-simple-server
wget https://github.com/aler9/rtsp-simple-server/releases/download/v0.21.0/rtsp-simple-server_v0.21.0_linux_amd64.tar.gz
tar -xf rtsp-simple-server_v0.21.0_linux_amd64.tar.gz
rm -f rtsp-simple-server_v0.21.0_linux_amd64.tar.gz
rm -f LICENSE
./rtsp-simple-server &
# record the live feed.
ffmpeg -i rtsp://127.0.0.1:8554/dji -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts saved_%Y-%m-%d_%H-%M-%S.ts
# upload the feed to aws
ffmpeg -re -stream_loop -1 -i $VIDEO_FILEPATH -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://$INGEST_ENDPOINT:443/app/$STREAM_KEY
```


# Do a anaylis of the live feed using cloud providers.

https://docs.aws.amazon.com/ivs/latest/userguide/getting-started-set-up-streaming.html

```

```
