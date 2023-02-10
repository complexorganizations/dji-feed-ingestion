#!/usr/bin/env bash

# Require script to be run as root
function super-user-check() {
    # This code checks to see if the script is running with root privileges.
    # If it is not, it will exit with an error message.
    if [ "${EUID}" -ne 0 ]; then
        echo "Error: You need to run this script as administrator."
        exit
    fi
}

# Check for root
super-user-check

# Get the current system information
function system-information() {
    # CURRENT_DISTRO is the ID of the current system
    # CURRENT_DISTRO_VERSION is the VERSION_ID of the current system
    if [ -f /etc/os-release ]; then
        # shellcheck source=/dev/null
        source /etc/os-release
        CURRENT_DISTRO=${ID}
        CURRENT_DISTRO_VERSION=${VERSION_ID}
    fi
}

# Get the current system information
system-information

# Pre-Checks system requirements
function installing-system-requirements() {
    if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ] || [ "${CURRENT_DISTRO}" == "fedora" ] || [ "${CURRENT_DISTRO}" == "centos" ] || [ "${CURRENT_DISTRO}" == "rhel" ] || [ "${CURRENT_DISTRO}" == "almalinux" ] || [ "${CURRENT_DISTRO}" == "rocky" ] || [ "${CURRENT_DISTRO}" == "arch" ] || [ "${CURRENT_DISTRO}" == "archarm" ] || [ "${CURRENT_DISTRO}" == "manjaro" ] || [ "${CURRENT_DISTRO}" == "alpine" ] || [ "${CURRENT_DISTRO}" == "freebsd" ] || [ "${CURRENT_DISTRO}" == "ol" ]; }; then
        if { [ ! -x "$(command -v git)" ] || [ ! -x "$(command -v curl)" ] || [ ! -x "$(command -v jq)" ] || [ ! -x "$(command -v ffmpeg)" ] || [ ! -x "$(command -v gst-launch-1.0)" ]; }; then
            if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ]; }; then
                apt-get update
                apt-get install pkg-config cmake m4 ffmpeg git procps build-essential jq libssl-dev libcurl4-openssl-dev liblog4cplus-dev libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-base-apps gstreamer1.0-plugins-bad gstreamer1.0-plugins-good gstreamer1.0-plugins-ugly gstreamer1.0-tools -y
            elif { [ "${CURRENT_DISTRO}" == "fedora" ] || [ "${CURRENT_DISTRO}" == "centos" ] || [ "${CURRENT_DISTRO}" == "rhel" ] || [ "${CURRENT_DISTRO}" == "almalinux" ] || [ "${CURRENT_DISTRO}" == "rocky" ]; }; then
                yum check-update
            elif { [ "${CURRENT_DISTRO}" == "arch" ] || [ "${CURRENT_DISTRO}" == "archarm" ] || [ "${CURRENT_DISTRO}" == "manjaro" ]; }; then
                pacman -Sy --noconfirm archlinux-keyring
            elif [ "${CURRENT_DISTRO}" == "alpine" ]; then
                apk update
            elif [ "${CURRENT_DISTRO}" == "freebsd" ]; then
                pkg update
            elif [ "${CURRENT_DISTRO}" == "ol" ]; then
                yum check-update
            fi
        fi
    else
        echo "Error: ${CURRENT_DISTRO} ${CURRENT_DISTRO_VERSION} is not supported."
        exit
    fi
}

# check for requirements
installing-system-requirements

# Global variables
AMAZON_KINESIS_VIDEO_STREAMS_GIT_PATH="https://github.com/awslabs/amazon-kinesis-video-streams-producer-sdk-cpp.git"
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH="/etc/amazon-kinesis-video-streams-producer-sdk-cpp"
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/build"
AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/open-source/local/lib"
AMAZON_KINESIS_VIDEO_STREAMS_PATH="./kvs_gstreamer_sample"
SYSTEM_IPV4=$(curl --ipv4 --connect-timeout 5 --tlsv1.3 --silent 'https://api.ipengine.dev' | jq -r '.network.ip')
# Create a key-value pair for the RTSP server and the Kinesis Video Streams Stream Name
declare -A RTSP_SERVERS
RTSP_SERVERS["rtsp://Administrator:Password@${SYSTEM_IPV4}:8554/drone_0"]="dji-stream-0"
RTSP_SERVERS["rtsp://Administrator:Password@${SYSTEM_IPV4}:8554/drone_1"]="dji-stream-1"
RTSP_SERVERS["rtsp://Administrator:Password@${SYSTEM_IPV4}:8554/drone_2"]="dji-stream-2"
RTSP_SERVERS["rtsp://Administrator:Password@${SYSTEM_IPV4}:8554/drone_3"]="dji-stream-3"
# AWS Credentials
AWS_ACCESS_KEY_ID="SAMPLEKEY"
AWS_SECRET_ACCESS_KEY="SAMPLESECRET"
AWS_DEFAULT_REGION="us-east-1"
# Kinesis Video Streams Bash Script
GITHUB_REPO_UPDATE_URL="https://raw.githubusercontent.com/complexorganizations/dji-feed-analysis/main/upload-feed-to-aws-kinesis.sh"
KINESIS_VIDEO_STREAMS_BASH_SERVICE="/etc/systemd/system/kinesis-video-streams-bash.service"
KINESIS_VIDEO_STREAMS_BASH_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/upload-feed-to-aws-kinesis.sh"
CURRENT_PATH_TO_SCRIPT=$(dirname "$(readlink -f "$0")")

# Build the application.
function build-kensis-application() {
    if [ ! -d "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}" ]; then
        git clone ${AMAZON_KINESIS_VIDEO_STREAMS_GIT_PATH} ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}
        mkdir -p ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        cmake -DBUILD_GSTREAMER_PLUGIN=TRUE -S ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH} -B ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        make -C ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        # Add the path to the .profile file so that it can be used in the future
        echo -e "export GST_PLUGIN_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}\nexport LD_LIBRARY_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH}:$LD_LIBRARY_PATH" >>~/.profile
        source ~/.profile
    fi
}

# Build the application.
build-kensis-application

# Install the script as a service.
function install-bash-as-service() {
    if [ ! -f ${KINESIS_VIDEO_STREAMS_BASH_PATH} ]; then
        # Note: Save the script in the correct directory.
        curl ${GITHUB_REPO_UPDATE_URL} -o ${KINESIS_VIDEO_STREAMS_BASH_PATH}
        chmod +x ${KINESIS_VIDEO_STREAMS_BASH_PATH}
    fi
    # Install the bash script as a service.
    if [ ! -f "${KINESIS_VIDEO_STREAMS_BASH_SERVICE}" ]; then
        echo "[Unit]
Wants=network.target
[Service]
ExecStart=${KINESIS_VIDEO_STREAMS_BASH_PATH}
[Install]
WantedBy=multi-user.target" >${KINESIS_VIDEO_STREAMS_BASH_SERVICE}
        if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
            systemctl daemon-reload
            systemctl enable kinesis-video-streams-bash
            systemctl restart kinesis-video-streams-bash
        elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
            service kinesis-video-streams-bash restart
        fi
    fi
}

# Install the bash script as a service.
install-bash-as-service

# Check the RTSP server status
function check-rtsp-server-status() {
    # Make sure the script is running in the correct directory
    if [ "${CURRENT_PATH_TO_SCRIPT}" == "${KINESIS_VIDEO_STREAMS_BASH_PATH}" ]; then
        # Create a loop that will run forever
        while true; do
            # Loop through the RTSP servers and check if they are alive
            for RTSP_SERVER in "${!RTSP_SERVERS[@]}"; do
                KINESIS_STREAM_NAME="${RTSP_SERVERS[${RTSP_SERVER}]}"
                RTSP_SERVER_LOG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/${KINESIS_STREAM_NAME}.log"
                # Check if a given RTSP server is alive and if it is than stream it
                if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER}" | wc -m)" -gt 100 ]; then
                    # Create a counter for the number of times the RTSP server has been started
                    RTSP_SERVER_START_COUNTER_${KINESIS_STREAM_NAME}=0
                    # Only start the RTSP stream to kinesis video streams if the counter is 0; this is to prevent the RTSP server from being started multiple times
                    if [ $RTSP_SERVER_START_COUNTER_${KINESIS_STREAM_NAME} == 0 ]; then
                        # We need to add a one to the counter so that the RTSP server is not started multiple times
                        RTSP_SERVER_START_COUNTER_${KINESIS_STREAM_NAME}=$((RTSP_SERVER_START_COUNTER_${KINESIS_STREAM_NAME} + 1))
                        # Start the RTSP server
                        AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM_NAME} "${RTSP_SERVER}" >${RTSP_SERVER_LOG} &
                        # Create a counter for the while loop
                        RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME} = 0
                        # Create a while loop that will check the health of the RTSP server
                        while [ $RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME} == 0 ]; do
                            # Check the status of the stream.
                            if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_LOG}" | wc -m)" -lt 100 ]; then
                                # End the stream to aws since the stream already eneded.
                                kill $!
                                RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME}=$((RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME} + 1))
                            fi
                            if [ "$(tail -n50 ${RTSP_SERVER_LOG} | grep 'Pad link failed' | wc -m)" -ge 1 ]; then
                                # End the stream if there is an issue
                                kill $!
                                RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME}=$((RTSP_SERVER_WHILE_COUNTER_${KINESIS_STREAM_NAME} + 1))
                            fi
                        done
                    fi
                fi
            done
        done
    fi
}

# Check if the RTSP server is alive and if it is than stream it
check-rtsp-server-status

# Note: This will not work i think since it seems to be broken in the logic of the app; i think i am not 100% sure but lets see what happens.
