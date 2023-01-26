#### Create a RTSP server.

### Creating the server and taking in the live stream
##### https://github.com/aler9/rtsp-simple-server/pull/1404

### Recording the live stream on the server
``` bash
ffmpeg -i rtsp://localhost:8554/test -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts /etc/rtsp-simple-server/main.ts
```

### Upload the feed from the server to aws.
``` bash
ffmpeg -re -stream_loop -1 -i /etc/rtsp-simple-server/main.ts -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://${INGEST_ENDPOINT}:443/app/${STREAM_KEY}
```

### Check the status of a stream.
```
ffprobe -v quiet -print_format json -show_streams rtmp://localhost:1935/test
```

#### Do a analysis of the live feed using cloud providers.

##### https://docs.aws.amazon.com/ivs/latest/userguide/getting-started-set-up-streaming.html

```

```

### Post Recording

##### What to do with the `.lrf` files? (DELETE ALL OF THEM)
``` bash
rm -f 100MEDIA/*.LRF
```

### Validate all the recorded videos are good.
``` bash
ffmpeg -v error -i first_input.mp4 -f null - 2 >> error.log
```

### Combine all the valid videos into a single video.
``` bash
ffmpeg -i concat:"first_input.mp4|second_input.mp4" output.mp4
```

### Combine all the srt file into one.
``` bash
ALL_FILES_IN_DIRECTORY=$(find 100MEDIA/  -type f -name "*.SRT")
for CURRENT_FILE_IN_DIRECTORY in ${ALL_FILES_IN_DIRECTORY}; do
  ALL_FILES_LIST[${ADD_CONTENT}]=${CURRENT_FILE_IN_DIRECTORY}
  ADD_CONTENT=$(("${ADD_CONTENT}" + 1))
done
for FILE_LIST in "${ALL_FILES_LIST[@]}"; do
  cat ${FILE_LIST} >> 100MEDIA/ALL.SRT
done
```
