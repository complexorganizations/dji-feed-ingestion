# Create a RTSP server.

1 . https://github.com/aler9/rtsp-simple-server


# Configure the stuff.

# Record the video live.
```
.\ffmpeg.exe -i "rtsp://127.0.0.1:8554/mystream" -c copy -f segment -strftime 1 -segment_time 60 -segment_format mpegts saved_%Y-%m-%d_%H-%M-%S.ts
```

# Feed the data into aws.
```
ffmpeg -re -stream_loop -1 -i $VIDEO_FILEPATH -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://$INGEST_ENDPOINT:443/app/$STREAM_KEY
```
https://docs.aws.amazon.com/ivs/latest/userguide/getting-started-set-up-streaming.html
