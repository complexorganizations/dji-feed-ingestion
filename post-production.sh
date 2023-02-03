# Global variables.
PATH_TO_MICRO_SD_CARD="100MEDIA"
PATH_TO_S3="s3://mybucket/myfolder"
RANDOM_DRONE_RUN_ID=$(openssl rand -hex 6)
DRONE_VIDEO_ZIP_NAME="${RANDOM_DRONE_RUN_ID}.zip"

# Remove all the temp .lrf files.
rm -f 100MEDIA/*.LRF

# Zip all the files.
zip -r ${DRONE_VIDEO_ZIP_NAME} ${PATH_TO_MICRO_SD_CARD}

# Upload all the files to a given service like S3
aws s3 cp ${DRONE_VIDEO_ZIP_NAME} ${PATH_TO_S3}

# Remove the content from the SD card.
rm -r ${PATH_TO_MICRO_SD_CARD}/*

# Using a serivce like lamda function unzip the files.
# Unzip all the videos and make sure all the video files are valid.
$ unzip ${DRONE_VIDEO_ZIP_NAME} -d ${RANDOM_DRONE_RUN_ID}

# Validate all the recordings are good.
ffmpeg -v error -i first_input.mp4 -f null - 2 >> error.log

# Combine all the recordings into one.
ffmpeg -i concat:"first_input.mp4|second_input.mp4" output.mp4

# Combine all the .srt files into one.
cat *.SRT >> all.srt

# analysis on the video depending on which is faster; the split videos or the whoel vidoes.
# Split the video frame by frame just the same way as aws; so we can convert that recorded data too.
