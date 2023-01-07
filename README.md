# Create a RTSP server.

1 . https://github.com/aler9/rtsp-simple-server


# Configure the stuff.



# Record the video live.
```
.\ffmpeg.exe -i "rtsp://127.0.0.1:8554/mystream" -c copy -f segment -strftime 1 -segment_time 60 -segment_format mpegts saved_%Y-%m-%d_%H-%M-%S.ts
```
