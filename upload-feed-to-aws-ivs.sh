# Global variables
RTSP_CONNECTION_STRING="rtsp://Administrator:Password@localhost:8554/test"
RTSP_VIDEO_PATH="/etc/rtsp-simple-server/main.ts"

# Always record the feed; even thought it might be offline since we don't know when it will come online.
ffmpeg -i ${RTSP_CONNECTION_STRING} -c copy -f segment -strftime 1 -segment_time 3600 -segment_format mpegts ${RTSP_VIDEO_PATH}

# Check the feed is online every minute; and if its online than feed it into aws.

# Prevent double feed into aws; only feed the data into aws once.

# Check if a given feed is alive.
if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_CONNECTION_STRING}" | wc -m)" -gt 100 ]; then
  echo "The feed is alive"
  ffmpeg -re -stream_loop -1 -i ${RTSP_VIDEO_PATH} -r 30 -c:v libx264 -pix_fmt yuv420p -profile:v main -preset veryfast -x264opts "nal-hrd=cbr:no-scenecut" -minrate 3000 -maxrate 3000 -g 60 -c:a aac -b:a 160k -ac 2 -ar 44100 -f flv rtmps://${INGEST_ENDPOINT}:443/app/${STREAM_KEY}
else
  exit
fi
