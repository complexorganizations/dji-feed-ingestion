#### Create a RTSP server.

### Creating the server and taking in the live stream
##### https://github.com/aler9/rtsp-simple-server/pull/1404

#### Check the status of a stream.
``` bash
ffprobe -v quiet -print_format json -show_streams rtsp://admin:password@localhost:8554/test
```

#### Watch the stream live.
``` bash
vlc rtsp://admin:password@localhost:8554/test
```

#### Recording the live stream on the server
``` bash
ffmpeg -i rtsp://localhost:8554/test -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts /etc/rtsp-simple-server/main.ts
```

#### Upload the feed from the server to aws.
``` bash
ffmpeg -re -stream_loop -1 -i /etc/rtsp-simple-server/main.ts -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://${INGEST_ENDPOINT}:443/app/${STREAM_KEY}
```

