# Global variables.
PATH_TO_MICRO_SD_CARD=$()
RANDOM_DRONE_RUN_ID=$(openssl rand -hex 6)
DRONE_VIDEO_ZIP_NAME="${RANDOM_DRONE_RUN_ID}.zip"

# Remove all the temp .lrf files.
rm -f 100MEDIA/*.LRF

# Validate all the recordings are good.
ffmpeg -v error -i first_input.mp4 -f null - 2 >> error.log

# Combine all the recordings into one.
ffmpeg -i concat:"first_input.mp4|second_input.mp4" output.mp4

# Zip all the files.
zip -r ${DRONE_VIDEO_ZIP_NAME} 100MEDIA/

# Upload all the files to a given service like S3

# Unzip all the videos and make sure all the video files are valid.
unzip ${DRONE_VIDEO_ZIP_NAME} -d ${RANDOM_DRONE_RUN_ID}

# Combine all the videos into a single video.

# Combine all the .srt files into one.

# Export the content out of here.

# Remove the content from the SD card.

# Upload the content to youtube.
