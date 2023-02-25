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
        if { [ ! -x "$(command -v cut)" ] || [ ! -x "$(command -v git)" ] || [ ! -x "$(command -v ffmpeg)" ]; }; then
            if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ]; }; then
                apt-get update
                apt-get install coreutils git ffmpeg curl openssl tar apt-transport-https ca-certificates gnupg -y
            elif { [ "${CURRENT_DISTRO}" == "ol" ] || [ "${CURRENT_DISTRO}" == "fedora" ] || [ "${CURRENT_DISTRO}" == "centos" ] || [ "${CURRENT_DISTRO}" == "rhel" ] || [ "${CURRENT_DISTRO}" == "almalinux" ] || [ "${CURRENT_DISTRO}" == "rocky" ]; }; then
                yum check-update
                yum install coreutils git ffmpeg curl openssl tar gpg -y
            elif { [ "${CURRENT_DISTRO}" == "arch" ] || [ "${CURRENT_DISTRO}" == "archarm" ] || [ "${CURRENT_DISTRO}" == "manjaro" ]; }; then
                pacman -Sy --noconfirm archlinux-keyring
            elif [ "${CURRENT_DISTRO}" == "alpine" ]; then
                apk update
            elif [ "${CURRENT_DISTRO}" == "freebsd" ]; then
                pkg update
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
RTSP_SIMPLE_SERVER_PATH="/etc/rtsp-simple-server"
RTSP_SIMPLE_SERVER_CONFIG="${RTSP_SIMPLE_SERVER_PATH}/rtsp-simple-server.yml"
RTSP_SIMPLE_SERVICE_APPLICATION="${RTSP_SIMPLE_SERVER_PATH}/rtsp-simple-server"
RTSP_SIMPLE_SERVER_SERVICE="/etc/systemd/system/rtsp-simple-server.service"
RTSP_CONFIG_FILE_GITHUB_URL="https://raw.githubusercontent.com/complexorganizations/dji-feed-analysis/main/assets/rtsp-simple-server.yml"
RTSP_SIMPLE_SERVER_LATEST_RELEASE=$(curl -s https://api.github.com/repos/aler9/rtsp-simple-server/releases/latest | grep browser_download_url | cut -d'"' -f4 | grep $(dpkg --print-architecture) | grep linux)
RTSP_SIMPLE_SERVER_LASTEST_FILE_NAME=$(echo "${RTSP_SIMPLE_SERVER_LATEST_RELEASE}" | cut --delimiter="/" --fields=9)
RTSP_SIMPLE_SERVER_TEMP_DOWNLOAD_PATH="/tmp/${RTSP_SIMPLE_SERVER_LASTEST_FILE_NAME}"

# Note: Get the latest release and don't get the "main" branch.
AMAZON_KINESIS_VIDEO_STREAMS_GIT_PATH="https://github.com/awslabs/amazon-kinesis-video-streams-producer-sdk-cpp.git"
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH="/etc/amazon-kinesis-video-streams-producer-sdk-cpp"
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/build"
AMAZON_KINESIS_VIDEO_STREAMS_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}/kvs_gstreamer_sample"
AMAZON_KINESIS_VIDEO_STREAMS_KVS_LOG_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/kvs_log_configuration"
AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/open-source/local/lib"
AMAZON_KINESIS_VIDEO_STREAMS_GST_STREAMER_CONFIG="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/src/gstreamer/gstkvssink.cpp"

CSP_CONNECTOR_PATH="/etc/csp-connector"
CSP_CONNECTOR_CONFIG="${CSP_CONNECTOR_PATH}/config.json"
CSP_CONNECTOR_SERVICE="/etc/systemd/system/csp-connector.service"
CSP_CONNECTOR_APPLICATION="${CSP_CONNECTOR_PATH}/csp-connector"
CSP_CONNECTOR_LATEST_RELEASE=$(curl -s https://api.github.com/repos/complexorganizations/csp-connector/releases/latest | grep browser_download_url | cut -d'"' -f4 | grep $(dpkg --print-architecture) | grep linux)
CSP_CONNECTOR_LATEST_FILE_NAME=$(echo "${CSP_CONNECTOR_LATEST_RELEASE}" | cut --delimiter="/" --fields=9)
CSP_CONNECTOR_TEMP_DOWNLOAD_PATH="/tmp/${CSP_CONNECTOR_LATEST_FILE_NAME}"

# Install rtsp application.
function install-rtsp-application() {
    mkdir -p ${RTSP_SIMPLE_SERVER_PATH}
    curl -L "${RTSP_SIMPLE_SERVER_LATEST_RELEASE}" -o ${RTSP_SIMPLE_SERVER_TEMP_DOWNLOAD_PATH}
    tar -xvf ${RTSP_SIMPLE_SERVER_TEMP_DOWNLOAD_PATH} -C ${RTSP_SIMPLE_SERVER_PATH}
    rm -f ${RTSP_SIMPLE_SERVER_TEMP_DOWNLOAD_PATH}
    curl ${RTSP_CONFIG_FILE_GITHUB_URL} -o ${RTSP_SIMPLE_SERVER_CONFIG}
    chmod +x ${RTSP_SIMPLE_SERVICE_APPLICATION}
    if [ ! -f "${RTSP_SIMPLE_SERVER_SERVICE}" ]; then
        # This code creates the service file
        # The service file is stored in /etc/systemd/system/rtsp-simple-server.service
        echo "[Unit]
Wants=network.target
[Service]
ExecStart=${RTSP_SIMPLE_SERVICE_APPLICATION} ${RTSP_SIMPLE_SERVER_CONFIG}
[Install]
WantedBy=multi-user.target" >${RTSP_SIMPLE_SERVER_SERVICE}
        if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
            systemctl daemon-reload
            systemctl enable rtsp-simple-server
            systemctl start rtsp-simple-server
        elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
            service rtsp-simple-server start
        fi
    fi
}

# Install the rtsp server.
install-rtsp-application

# Build the application.
function build-kensis-application() {
    if [ ! -d "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}" ]; then
        if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ]; }; then
            apt-get install build-essential pkg-config cmake m4 libssl-dev libcurl4-openssl-dev liblog4cplus-dev libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-base-apps gstreamer1.0-plugins-bad gstreamer1.0-plugins-good gstreamer1.0-plugins-ugly gstreamer1.0-tools -y
        elif { [ "${CURRENT_DISTRO}" == "ol" ] || [ "${CURRENT_DISTRO}" == "fedora" ] || [ "${CURRENT_DISTRO}" == "centos" ] || [ "${CURRENT_DISTRO}" == "rhel" ] || [ "${CURRENT_DISTRO}" == "almalinux" ] || [ "${CURRENT_DISTRO}" == "rocky" ]; }; then
            yum install gcc gcc-c++ make cmake openssl-devel libcurl-devel log4cplus-devel gstreamer1-devel gstreamer1-plugins-base-devel gstreamer1-plugins-good gstreamer1-plugins-bad-free gstreamer1-plugins-bad-free-extras gstreamer1-plugins-bad-freeworld gstreamer1-plugins-bad-nonfree gstreamer1-plugins-ugly gstreamer1-plugins-ugly-free gstreamer1-plugins-ugly-free-extras gstreamer1-plugins-ugly-freeworld gstreamer1-plugins-ugly-nonfree -y
        elif { [ "${CURRENT_DISTRO}" == "arch" ] || [ "${CURRENT_DISTRO}" == "archarm" ] || [ "${CURRENT_DISTRO}" == "manjaro" ]; }; then
            pacman -S base-devel cmake openssl curl log4cplus gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly gst-libav -y
        elif [ "${CURRENT_DISTRO}" == "alpine" ]; then
            apk add build-base cmake openssl-dev curl-dev log4cplus-dev gstreamer-dev gstreamer-plugins-base-dev gstreamer-plugins-good gstreamer-plugins-bad gstreamer-plugins-ugly gstreamer-ffmpeg -y
        elif [ "${CURRENT_DISTRO}" == "freebsd" ]; then
            pkg install cmake openssl curl log4cplus gstreamer1-plugins-base gstreamer1-plugins-good gstreamer1-plugins-bad gstreamer1-plugins-ugly gstreamer1-plugins-libav -y
        fi
        git clone ${AMAZON_KINESIS_VIDEO_STREAMS_GIT_PATH} ${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}
        # Change the path to the log file so the correct path is build everytime.
        # sed -i "s|../kvs_log_configuration|${AMAZON_KINESIS_VIDEO_STREAMS_KVS_LOG_PATH}|g" ${AMAZON_KINESIS_VIDEO_STREAMS_GST_STREAMER_CONFIG}
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

# Run the application.
# AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} ${AMAZON_KINESIS_VIDEO_STREAMS_PATH} ${KINESIS_STREAM} "${RTSP_SERVER}"

# Install Google Cloud
function install-google-cloud() {
    if { [ ! -x "$(command -v gcloud)" ] || [ ! -x "$(command -v vaictl)" ]; }; then
        if { [ "${CURRENT_DISTRO}" == "ubuntu" ] || [ "${CURRENT_DISTRO}" == "debian" ] || [ "${CURRENT_DISTRO}" == "raspbian" ] || [ "${CURRENT_DISTRO}" == "pop" ] || [ "${CURRENT_DISTRO}" == "kali" ] || [ "${CURRENT_DISTRO}" == "linuxmint" ] || [ "${CURRENT_DISTRO}" == "neon" ]; }; then
            echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
            curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
            apt-get update
            apt-get install google-cloud-cli -y
            # gcloud auth login --no-launch-browser
            # gcloud auth application-default login --no-launch-browser
            # gcloud services enable visionai.googleapis.com
            # Install Google cloud vision ai
            curl -L https://github.com/google/visionai/releases/download/v0.0.4/visionai_0.0-4_amd64.deb -o visionai_0.0-4_amd64.deb
            apt-get install ./visionai_0.0-4_amd64.deb
            rm -f visionai_0.0-4_amd64.deb
        elif { [ "${CURRENT_DISTRO}" == "ol" ] || [ "${CURRENT_DISTRO}" == "fedora" ] || [ "${CURRENT_DISTRO}" == "centos" ] || [ "${CURRENT_DISTRO}" == "rhel" ] || [ "${CURRENT_DISTRO}" == "almalinux" ] || [ "${CURRENT_DISTRO}" == "rocky" ]; }; then
            yum check-update
        elif { [ "${CURRENT_DISTRO}" == "arch" ] || [ "${CURRENT_DISTRO}" == "archarm" ] || [ "${CURRENT_DISTRO}" == "manjaro" ]; }; then
            pacman -Sy
        elif [ "${CURRENT_DISTRO}" == "alpine" ]; then
            apk update
        elif [ "${CURRENT_DISTRO}" == "freebsd" ]; then
             pkg update
        fi
    fi
}

# Install Google Cloud
install-google-cloud

# Feed the data into google cloud vision ai
# vaictl -p github-code-snippets -l us-central1 -c application-cluster-0 --service-endpoint visionai.googleapis.com send rtsp to streams dji-stream-0 --rtsp-uri rtsp://Administrator:Password@localhost:8554/drone_0

# Install the cloud connector.
function install-cps-connetor() {
    if [ ! -d "${CSP_CONNECTOR_PATH}" ]; then
        mkdir -p ${CSP_CONNECTOR_PATH}
        curl -L "${CSP_CONNECTOR_LATEST_RELEASE}" -o ${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}
        tar -xvf ${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH} -C ${CSP_CONNECTOR_PATH}
        rm -f ${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}
        chmod +x ${CSP_CONNECTOR_APPLICATION}
        if [ ! -f "${CSP_CONNECTOR_SERVICE}" ]; then
            # This code creates the service file
            # The service file is stored in /etc/systemd/system/csp-connector.service
            echo "[Unit]
Wants=network.target
[Service]
ExecStart=${CSP_CONNECTOR_APPLICATION} -config=${CSP_CONNECTOR_CONFIG}
[Install]
WantedBy=multi-user.target" >${CSP_CONNECTOR_SERVICE}
            if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
                systemctl daemon-reload
                systemctl enable csp-connector
                systemctl start csp-connector
            elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
                service csp-connector start
            fi
        fi
    fi
}

# Install the cloud connector
# install-cps-connetor
