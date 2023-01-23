# Create a RTSP server.

##### https://github.com/aler9/rtsp-simple-server/pull/1404
```
apt-get update
apt-get upgrade -y
apt-get install ffmpeg wget coreutils screen -y
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

```
sed "s|runOnDemand:|runOnDemand: ffmpeg -re -stream_loop -1 -i main.tf -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://4e99170d7608.global-contribute.live-video.net:443/app/sk_us-east-1_2Je2yx6VbfSk_If4lY1OevQwx8gQ8JnZQi1qOlbGWIj|g' /etc/rtsp-simple-server/rtsp-simple-server.yml
sed 's|runOnReady:|runOnReady: ffmpeg -i rtsp://127.0.0.1:8554/dji -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts main.tf|g' /etc/rtsp-simple-server/rtsp-simple-server.yml
```

# Do a analysis of the live feed using cloud providers.

##### https://docs.aws.amazon.com/ivs/latest/userguide/getting-started-set-up-streaming.html

```

```

# Validate all the recorded videos are good.
``` bash
ALL_FILES_IN_DIRECTORY=$(find PATH_TO_VIDEO_FILES/ -type f -name "*.MP4")
for CURRENT_FILE_IN_DIRECTORY in ${ALL_FILES_IN_DIRECTORY}; do
  ALL_FILES_LIST[${ADD_CONTENT}]=${CURRENT_FILE_IN_DIRECTORY}
  ADD_CONTENT=$(("${ADD_CONTENT}" + 1))
done
for FILE_LIST in "${ALL_FILES_LIST[@]}"; do
  ffmpeg -v error -i ${FILE_LIST} -f null - 2>error.log
done
```

# Combine all the videos into a single video.
```
ffmpeg -i concat:"input1.mp4|input2.mp4" output.mp4
```
