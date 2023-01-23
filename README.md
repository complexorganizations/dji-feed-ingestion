#### Create a RTSP server.

### Creating the server and taking in the live stream
##### https://github.com/aler9/rtsp-simple-server/pull/1404

### Recording the live stream on the server
``` bash
ffmpeg -i rtsp://127.0.0.1:8554/dji -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts saved_%Y-%m-%d_%H-%M-%S.ts
```

### Upload the feed from the server to aws.
``` bash
ffmpeg -re -stream_loop -1 -i ${VIDEO_FILEPATH} -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://${INGEST_ENDPOINT}:443/app/${STREAM_KEY}
```

#### Do a analysis of the live feed using cloud providers.

##### https://docs.aws.amazon.com/ivs/latest/userguide/getting-started-set-up-streaming.html

```

```

### Post Recording
### Validate all the recorded videos are good.
``` bash
ffmpeg -v error -i first_input.mp4 -f null - 2 >> error.log
```

### Combine all the valid videos into a single video.
``` bash
ffmpeg -i concat:"first_input.mp4|second_input.mp4" output.mp4
```
