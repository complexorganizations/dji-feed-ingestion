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
        if { [ ! -x "$(command -v git)" ] || [ ! -x "$(command -v curl)" ] || [ ! -x "$(command -v jq)" ]; }; then
            if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ]; }; then
                apt-get update
                apt-get install pkg-config cmake m4 git procps build-essential jq libssl-dev libcurl4-openssl-dev liblog4cplus-dev libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-base-apps gstreamer1.0-plugins-bad gstreamer1.0-plugins-good gstreamer1.0-plugins-ugly gstreamer1.0-tools -y
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
AMAZON_KINESIS_VIDEO_STREAMS_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/kvs_gstreamer_sample"
SYSTEM_IPV4=$(curl --ipv4 --connect-timeout 5 --tlsv1.3 --silent 'https://api.ipengine.dev' | jq -r '.network.ip')
# RTSP Paths
RTSP_SERVER_ZERO="rtsp://admin:password@${SYSTEM_IPV4}:8554/zero"
RTSP_SERVER_ONE="rtsp://admin:password@${SYSTEM_IPV4}:8554/one"
RTSP_SERVER_TWO="rtsp://admin:password@${SYSTEM_IPV4}:8554/two"
RTSP_SERVER_THREE="rtsp://admin:password@${SYSTEM_IPV4}:8554/three"
# Kinesis Video Streams Variables
KINESIS_STREAM_ZERO="rtsp-stream-0"
KINESIS_STREAM_ONE="rtsp-stream-1"
KINESIS_STREAM_TWO="rtsp-stream-2"
KINESIS_STREAM_THREE="rtsp-stream-3"
# AWS Credentials
AWS_ACCESS_KEY_ID="SAMPLEKEY"
AWS_SECRET_ACCESS_KEY="SAMPLESECRET"
AWS_DEFAULT_REGION="us-east-1"
# Logging
RTSP_SERVER_ZERO_LOG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/rtsp-server-zero.log"
RTSP_SERVER_ONE_LOG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/rtsp-server-one.log"
RTSP_SERVER_TWO_LOG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/rtsp-server-two.log"
RTSP_SERVER_THREE_LOG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/rtsp-server-three.log"

# Build the application.
function build-kensis-application() {
    if [ ! -d "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}" ]; then
        git clone ${AMAZON_KINESIS_VIDEO_STREAMS_GIT_PATH} ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}
        mkdir -p ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        cmake -DBUILD_GSTREAMER_PLUGIN=TRUE -S ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH} -B ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        make -C ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}
        # Add the path to the .profile file so that it can be used in the future
        echo -e "export GST_PLUGIN_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}\nexport LD_LIBRARY_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH}" >>~/.profile
        source ~/.profile
    fi
}

# Build the application.
build-kensis-application

# Check the RTSP server status
function check-rtsp-server-status() {
    while true; do
        # Loop through the RTSP servers and check if they are alive
        # Check if a given RTSP server is alive and if it is than stream it
        # Only run the stream once.
        if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_ZERO}" | wc -m)" -gt 100 ]; then
            # Counter for the while loop
            RTSP_SERVER_ZERO_COUNTER=0
            if [ ${RTSP_SERVER_ZERO_COUNTER} == 0 ]; then
                # Add 1 to start the loop.
                RTSP_SERVER_ZERO_COUNTER=$((RTSP_SERVER_ZERO_COUNTER + 1))
                # Start kensis
                AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ./${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM_ZERO} "${RTSP_SERVER_ZERO}" > ${RTSP_SERVER_ZERO_LOG} &
                # Counter for the while loop.
                RTSP_SERVER_ZERO_CHECK_COUNTER=0
                while [ ${RTSP_SERVER_ZERO_CHECK_COUNTER} -le 0 ]; do
                    # Check the status of the stream.
                    if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_ZERO}" | wc -m)" -lt 100 ]; then
                        # End the stream to aws since the stream already eneded.
                        kill $!
                        RTSP_SERVER_ZERO_CHECK_COUNTER=$((RTSP_SERVER_ZERO_CHECK_COUNTER + 1))
                    fi
                    if [ "$(tail -n50 ${RTSP_SERVER_ZERO_LOG} | grep 'pad link failed' | wc -m)" -ge 1 ]; then
                        # End the stream if there is an issue
                        kill $!
                        RTSP_SERVER_ZERO_CHECK_COUNTER=$((RTSP_SERVER_ZERO_CHECK_COUNTER + 1))
                    fi
                    sleep 15
                done
                RTSP_SERVER_ZERO_COUNTER=$((RTSP_SERVER_ZERO_COUNTER - 1))
            fi
        fi
        if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_ONE}" | wc -m)" -gt 100 ]; then
            # Counter for the while loop
            RTSP_SERVER_ONE_COUNTER=0
            if [ ${RTSP_SERVER_ONE_COUNTER} == 0 ]; then
                # Add 1 to start the loop.
                RTSP_SERVER_ONE_COUNTER=$((RTSP_SERVER_ONE_COUNTER + 1))
                # Start kensis
                AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ./${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM_ONE} "${RTSP_SERVER_ONE}" > ${RTSP_SERVER_ONE_LOG} &
                # Counter for the while loop.
                RTSP_SERVER_ONE_CHECK_COUNTER=0
                while [ ${RTSP_SERVER_ONE_CHECK_COUNTER} -le 0 ]; do
                    # Check the status of the stream.
                    if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_ONE}" | wc -m)" -lt 100 ]; then
                        # End the stream to aws since the stream already eneded.
                        kill $!
                        RTSP_SERVER_ONE_CHECK_COUNTER=$((RTSP_SERVER_ONE_CHECK_COUNTER + 1))
                    fi
                    if [ "$(tail -n50 ${RTSP_SERVER_ONE_LOG} | grep 'pad link failed' | wc -m)" -ge 1 ]; then
                        # End the stream if there is an issue
                        kill $!
                        RTSP_SERVER_ONE_CHECK_COUNTER=$((RTSP_SERVER_ONE_CHECK_COUNTER + 1))
                    fi
                    sleep 15
                done
                RTSP_SERVER_ONE_COUNTER=$((RTSP_SERVER_ONE_COUNTER - 1))
            fi
        fi
        if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_TWO}" | wc -m)" -gt 100 ]; then
            # Counter for the while loop
            RTSP_SERVER_TWO_COUNTER=0
            if [ ${RTSP_SERVER_TWO_COUNTER} == 0 ]; then
                # Add 1 to start the loop.
                RTSP_SERVER_TWO_COUNTER=$((RTSP_SERVER_TWO_COUNTER + 1))
                # Start kensis
                AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ./${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM_TWO} "${RTSP_SERVER_TWO}" > ${RTSP_SERVER_TWO_LOG} &
                # Counter for the while loop.
                RTSP_SERVER_TWO_CHECK_COUNTER=0
                while [ ${RTSP_SERVER_TWO_CHECK_COUNTER} -le 0 ]; do
                    # Check the status of the stream.
                    if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_TWO}" | wc -m)" -lt 100 ]; then
                        # End the stream to aws since the stream already eneded.
                        kill $!
                        RTSP_SERVER_TWO_CHECK_COUNTER=$((RTSP_SERVER_TWO_CHECK_COUNTER + 1))
                    fi
                    if [ "$(tail -n50 ${RTSP_SERVER_TWO_LOG} | grep 'pad link failed' | wc -m)" -ge 1 ]; then
                        # End the stream if there is an issue
                        kill $!
                        RTSP_SERVER_TWO_CHECK_COUNTER=$((RTSP_SERVER_TWO_CHECK_COUNTER + 1))
                    fi
                    sleep 15
                done
                RTSP_SERVER_TWO_COUNTER=$((RTSP_SERVER_TWO_COUNTER - 1))
            fi
        fi
        if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_THREE}" | wc -m)" -gt 100 ]; then
            # Counter for the while loop
            RTSP_SERVER_THREE_COUNTER=0
            if [ ${RTSP_SERVER_THREE_COUNTER} == 0 ]; then
                # Add 1 to start the loop.
                RTSP_SERVER_THREE_COUNTER=$((RTSP_SERVER_THREE_COUNTER + 1))
                # Start kensis
                AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ./${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM_THREE} "${RTSP_SERVER_THREE}" > ${RTSP_SERVER_THREE_LOG} &
                # Counter for the while loop.
                RTSP_SERVER_THREE_CHECK_COUNTER=0
                while [ ${RTSP_SERVER_THREE_CHECK_COUNTER} -le 0 ]; do
                    # Check the status of the stream.
                    if [ "$(ffprobe -v quiet -print_format json -show_streams "${RTSP_SERVER_THREE}" | wc -m)" -lt 100 ]; then
                        # End the stream to aws since the stream already eneded.
                        kill $!
                        RTSP_SERVER_THREE_CHECK_COUNTER=$((RTSP_SERVER_THREE_CHECK_COUNTER + 1))
                    fi
                    if [ "$(tail -n50 ${RTSP_SERVER_THREE_LOG} | grep 'pad link failed' | wc -m)" -ge 1 ]; then
                        # End the stream if there is an issue
                        kill $!
                        RTSP_SERVER_THREE_CHECK_COUNTER=$((RTSP_SERVER_THREE_CHECK_COUNTER + 1))
                    fi
                    sleep 15
                done
                RTSP_SERVER_THREE_COUNTER=$((RTSP_SERVER_THREE_COUNTER - 1))
            fi
        fi
        sleep 15
    done
}

# Check if the RTSP server is alive and if it is than stream it
check-rtsp-server-status
