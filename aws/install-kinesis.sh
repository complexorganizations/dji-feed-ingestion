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
