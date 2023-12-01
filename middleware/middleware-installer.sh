#!/usr/bin/env bash
# https://github.com/complexorganizations/dji-feed-ingestion

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
        CURRENT_DISTRO_VERSION_MAJOR=$(echo "${CURRENT_DISTRO_VERSION}" | cut --delimiter="." --fields=1)
    fi
}

# Get the current system information
system-information

# Pre-Checks system requirements
function installing-system-requirements() {
    if { [ "${CURRENT_DISTRO}" == "ubuntu" ] && [ "${CURRENT_DISTRO_VERSION_MAJOR}" == "22" ]; }; then
        if { [ ! -x "$(command -v cut)" ] || [ ! -x "$(command -v git)" ] || [ ! -x "$(command -v ffmpeg)" ] || [ ! -x "$(command -v zip)" ] || [ ! -x "$(command -v unzip)" ] || [ ! -x "$(command -v systemd-detect-virt)" ] || [ ! -x "$(command -v python)" ]; }; then
            if [ "${CURRENT_DISTRO}" == "ubuntu" ]; then
                apt-get update
                apt-get update && apt-get upgrade -y && apt-get dist-upgrade -y && apt-get clean -y && apt-get autoremove -y && apt-get autoclean -y && apt-get install -f -y
                apt-get install coreutils git ffmpeg curl openssl tar apt-transport-https ca-certificates gnupg zip unzip systemd python3 python3-pip -y
                # Install the requirements for the pip packages
                if [ -x "$(command -v pip)" ]; then
                    pip install cmake
                fi
            fi
        fi
    else
        echo "Error: ${CURRENT_DISTRO}, ${CURRENT_DISTRO_VERSION} is not supported."
        exit
    fi
}

# check for requirements
installing-system-requirements

# Checking For Virtualization
function virt-check() {
    # This code checks if the system is running in a supported virtualization.
    # It returns the name of the virtualization if it is supported, or "none" if
    # it is not supported. This code is used to check if the system is running in
    # a virtual machine, and if so, if it is running in a supported virtualization.
    CURRENT_SYSTEM_VIRTUALIZATION=$(systemd-detect-virt --container)
    case ${CURRENT_SYSTEM_VIRTUALIZATION} in
    "docker" | "none" | "wsl") ;;
    *)
        echo "${CURRENT_SYSTEM_VIRTUALIZATION} virtualization is not supported (yet)."
        exit
        ;;
    esac
}

# Virtualization Check
virt-check

# Make sure the script is running inside docker or else exit the script.
function check-inside-docker() {
    if [ ! -f /.dockerenv ]; then
        echo "Error: This script isn't running inside docker."
        # exit
        # Note: Remove the comment above and close the app in production.
    fi
}

# Make sure the application is running inside docker.
check-inside-docker

# The following function checks if the current init system is one of the allowed options.
function check-current-init-system() {
    # This function checks if the current init system is systemd or sysvinit.
    # If it is neither, the script exits.
    CURRENT_INIT_SYSTEM=$(ps --no-headers -o comm 1)
    # This line retrieves the current init system by checking the process name of PID 1.
    case ${CURRENT_INIT_SYSTEM} in
    # The case statement checks if the retrieved init system is one of the allowed options.
    *"systemd"* | *"init"*)
        # If the init system is systemd or sysvinit (init), continue with the script.
        ;;
    *)
        # If the init system is not one of the allowed options, display an error message and exit.
        echo "${CURRENT_INIT_SYSTEM} init is not supported (yet)."
        exit
        ;;
    esac
}

# The check-current-init-system function is being called.
check-current-init-system
# Calls the check-current-init-system function.

# Global variables
# Assigns the latest release of MediaMTX to a variable
MEDIAMTX_LATEST_RELEASE=$(curl -s https://api.github.com/repos/bluenviron/mediamtx/releases/latest | grep browser_download_url | cut --delimiter='"' --fields=4 | grep "$(dpkg --print-architecture)" | grep linux)
# Extracts the file name from the latest release URL and assigns it to a variable
MEDIAMTX_LASTEST_FILE_NAME=$(echo "${MEDIAMTX_LATEST_RELEASE}" | cut --delimiter="/" --fields=9)
# Assigns a temporary download path for the MediaMTX zip file
MEDIAMTX_TEMP_DOWNLOAD_PATH="/tmp/${MEDIAMTX_LASTEST_FILE_NAME}"
# Assigns a URL for the MediaMTX configuration file
MEDIAMTX_CONFIG_FILE_GITHUB_URL="https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/middleware/mediamtx.yml"
# Assigns a path for the MediaMTX directory
MEDIAMTX_LOCAL_PATH="/etc/mediamtx"
# Assigns a path for the mediamtx configuration file
MEDIAMTX_LOCAL_CONFIG_PATH="${MEDIAMTX_LOCAL_PATH}/mediamtx.yml"
# Assigns a path for the mediamtx service file
MEDIAMTX_SERVICE_FILE_PATH="/etc/systemd/system/mediamtx.service"
# Assigns a path for the mediamtx binary
MEDIAMTX_BINARY_PATH="${MEDIAMTX_LOCAL_PATH}/mediamtx"
# The variable to stream a test video feed as an test connection.
MEDIAMTX_TEST_CONNECTION_ZERO="rtsp://PublishAdministrator:PublishPassword@localhost:8554/test_zero"
# The path in the system that will host the test feed.
MEDIAMTX_TEST_FEED_ZERO_SERVICE_PATH="/etc/systemd/system/mediamtx-test-feed-zero.service"
# The path to the video file where the video is hosted.
MEDIAMTX_TEST_VIDEO_PATH_ZERO="${MEDIAMTX_LOCAL_PATH}/output_zero.mp4"

# Assigns the latest release of the Amazon Kinesis Video Streams Producer SDK to a variable
AMAZON_KINESIS_VIDEO_STREAMS_LATEST_RELEASE=$(curl -s https://api.github.com/repos/awslabs/amazon-kinesis-video-streams-producer-sdk-cpp/releases/latest | grep zipball_url | cut -d'"' -f4)
# Extracts the file name from the latest release URL and assigns it to a variable
AMAZON_KINESIS_VIDEO_STREAMS_FILE_NAME=$(echo "${AMAZON_KINESIS_VIDEO_STREAMS_LATEST_RELEASE}" | cut --delimiter="/" --fields=6)
# Assigns a path for the Kinesis Video Streams Producer SDK
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH="/etc/${AMAZON_KINESIS_VIDEO_STREAMS_FILE_NAME}"
# Assigns a path for the Kinesis Video Streams Producer SDK local libraries
AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/open-source/local/lib"
# Assigns a path for building the Kinesis Video Streams Producer SDK
AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH="${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}/build"
# Assigns a temporary download path for the Kinesis Video Streams Producer SDK zip file
AMAZON_KINESIS_VIDEO_STREAMS_TEMP_DOWNLOAD_PATH="/tmp/${AMAZON_KINESIS_VIDEO_STREAMS_FILE_NAME}.zip"

# Assigns the latest release of the AWS CLI to a variable
AMAZON_CLI_LATEST_RELEASE="https://awscli.amazonaws.com/awscli-exe-linux-$(arch).zip"
# Extracts the file name from the latest release URL and assigns it to a variable
AMAZON_CLI_FILE_NAME=$(echo "${AMAZON_CLI_LATEST_RELEASE}" | cut --delimiter="/" --fields=4)
# Assigns a temporary download path for the AWS CLI zip file
AMAZON_CLI_TEMP_DOWNLOAD_PATH="/tmp/${AMAZON_CLI_FILE_NAME}"
# Assigns a temporary path for the AWS CLI installation files to be extracted to before installation.
AMAZON_CLI_TEMP_INSTALL_PATH="/tmp/aws/"
# The path to the aws installation script
AMAZON_CLI_INSTALL_SCRIPT_PATH="${AMAZON_CLI_TEMP_INSTALL_PATH}install"

# Assigns the latest release of the CSP Connector to a variable
CSP_CONNECTOR_LATEST_RELEASE=$(curl -s https://api.github.com/repos/complexorganizations/dji-feed-ingestion/releases/latest | grep browser_download_url | cut --delimiter='"' --fields=4 | grep "$(dpkg --print-architecture)" | grep linux)
# Assigns the config file for the CSP connector.
CSP_CONNECTOR_CONFIG_URL="https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/middleware/config.json"
# Assigns a path for the CSP Connector
CSP_CONNECTOR_PATH="/etc/csp-connector"
# Assigns a path for the CSP Connector configuration file
CSP_CONNECTOR_CONFIG="${CSP_CONNECTOR_PATH}/config.json"
# Assigns a path for the CSP Connector log file.
CSP_CONNECTOR_LOG_FILE="${CSP_CONNECTOR_PATH}/log.txt"
# Extracts the file name from the latest release URL and assigns it to a variable
CSP_CONNECTOR_LATEST_FILE_NAME=$(echo "${CSP_CONNECTOR_LATEST_RELEASE}" | cut --delimiter="/" --fields=9 | cut --delimiter="-" --fields=1-2)
# Assigns a path for the CSP Connector application
CSP_CONNECTOR_APPLICATION="${CSP_CONNECTOR_PATH}/${CSP_CONNECTOR_LATEST_FILE_NAME}"
# Assigns a path for the CSP Connector service file
CSP_CONNECTOR_SERVICE="/etc/systemd/system/csp-connector.service"
# Assigns a temporary download path for the CSP Connector zip file
CSP_CONNECTOR_TEMP_DOWNLOAD_PATH="/tmp/${CSP_CONNECTOR_LATEST_FILE_NAME}"
# Assigns a permanent download path for the CSP Connector binary file
CSP_CONNECTOR_BINARY_PATH="/usr/bin/csp-connector"

# Assigns the latest release of the Google Cloud Vision AI to a variable
GOOGLE_CLOUD_VISION_AI_LATEST_RELEASE=$(curl -s https://api.github.com/repos/google/visionai/releases/latest | grep browser_download_url | cut --delimiter='"' --fields=4 | grep ".deb")
# Extracts the file name from the latest release URL and assigns it to a variable
GOOGLE_CLOUD_VISION_AI_LEAST_FILE_NAME=$(echo "${GOOGLE_CLOUD_VISION_AI_LATEST_RELEASE}" | cut --delimiter="/" --fields=9 | grep ".deb")
# Assigns a temporary download path for the Google Cloud Vision AI zip file
GOOGLE_CLOUD_VISION_AI_TEMP_DOWNLOAD_PATH="/tmp/${GOOGLE_CLOUD_VISION_AI_LEAST_FILE_NAME}"

# Get the latest release of youtube DLP
YOUTUBE_DLP_LATEST_RELEASE_URL=$(curl -s https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest | grep browser_download_url | cut --delimiter='"' --fields=4 | grep linux | head -n1)
# The system's local path where the yt-dlp should be placed
YOUTUBE_DLP_LOCAL_PATH="/usr/bin/yt-dlp"
# Test video to download and evaluate from YouTube
YOUTUBE_DLP_TEST_VIDEO_URL_ZERO="https://www.youtube.com/watch?v=IQ3SaSf--8Q"

# Install mediamtx application.
function install-mediamtx-application() {
    if [ ! -d "${MEDIAMTX_LOCAL_PATH}" ]; then
        # Create the directory.
        mkdir -p "${MEDIAMTX_LOCAL_PATH}"
        # Download the application.
        curl -L "${MEDIAMTX_LATEST_RELEASE}" -o "${MEDIAMTX_TEMP_DOWNLOAD_PATH}"
        # Extract the application.
        tar -xvf "${MEDIAMTX_TEMP_DOWNLOAD_PATH}" -C "${MEDIAMTX_LOCAL_PATH}"
        # Remove the downloaded file.
        rm -f "${MEDIAMTX_TEMP_DOWNLOAD_PATH}"
        # Download the configuration file.
        curl -L "${MEDIAMTX_CONFIG_FILE_GITHUB_URL}" -o "${MEDIAMTX_LOCAL_CONFIG_PATH}"
        # Change the permissions.
        chmod +x ${MEDIAMTX_BINARY_PATH}
        if [ ! -f "${MEDIAMTX_SERVICE_FILE_PATH}" ]; then
            # This code creates the service file
            # The service file is stored in /etc/systemd/system/mediamtx.service
            echo "[Unit]
Wants=network.target
[Service]
ExecStart=${MEDIAMTX_BINARY_PATH} ${MEDIAMTX_LOCAL_CONFIG_PATH}
[Install]
WantedBy=multi-user.target" >${MEDIAMTX_SERVICE_FILE_PATH}
            # Reload the daemon.
            systemctl daemon-reload
            if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
                systemctl enable --now mediamtx
                systemctl start mediamtx
            elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
                service mediamtx start
            fi
        fi
    fi
}

# Install the mediamtx server.
install-mediamtx-application

# Setup the rest feed for the stream
function setup-test-feed() {
    # Check if youtube dlp is installed
    if [ ! -x "$(command -v yt-dlp)" ]; then
        # Install youtube dlp
        curl -L "${YOUTUBE_DLP_LATEST_RELEASE_URL}" -o ${YOUTUBE_DLP_LOCAL_PATH}
        chmod +x ${YOUTUBE_DLP_LOCAL_PATH}
    fi
    # Downloading the zero test video.
    if [ ! -f ${MEDIAMTX_TEST_VIDEO_PATH_ZERO} ]; then
        yt-dlp -S ext:mp4:m4a "${YOUTUBE_DLP_TEST_VIDEO_URL_ZERO}" -o ${MEDIAMTX_TEST_VIDEO_PATH_ZERO}
    fi
    # Create a test feed if it does not exist already
    if [ ! -f "${MEDIAMTX_TEST_FEED_ZERO_SERVICE_PATH}" ]; then
        echo "[Unit]
Description=Create a test feed for MediaMTX server #0
Wants=network.target
[Service]
Type=simple
ExecStart=ffmpeg -re -stream_loop -1 -i ${MEDIAMTX_TEST_VIDEO_PATH_ZERO} -c copy -f rtsp ${MEDIAMTX_TEST_CONNECTION_ZERO}
Restart=always
[Install]
WantedBy=multi-user.target" >${MEDIAMTX_TEST_FEED_ZERO_SERVICE_PATH}
    fi
    # Reload the daemon
    systemctl daemon-reload
    if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
        systemctl enable --now mediamtx-test-feed-zero
        systemctl start mediamtx-test-feed-zero
    elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
        service mediamtx-test-feed-zero start
    fi
}

# Setup the test feed
setup-test-feed

# Build the application.
function build-kensis-application() {
    if [ ! -d "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}" ]; then
        # Install the dependencies.
        apt-get install libssl-dev libcurl4-openssl-dev libunwind-dev libgstreamer1.0-dev liblog4cplus-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-base-apps gstreamer1.0-plugins-bad gstreamer1.0-plugins-good gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools build-essential pkg-config cmake m4 byacc curl g++ git maven openjdk-8-jdk python2.7 -y
        # Download the application.
        curl -L "${AMAZON_KINESIS_VIDEO_STREAMS_LATEST_RELEASE}" -o "${AMAZON_KINESIS_VIDEO_STREAMS_TEMP_DOWNLOAD_PATH}"
        # Extract the application.
        unzip "${AMAZON_KINESIS_VIDEO_STREAMS_TEMP_DOWNLOAD_PATH}" -d /etc/
        # Remove the downloaded file.
        mv /etc/awslabs-amazon-kinesis-video-streams-producer-sdk-cpp-* /etc/"${AMAZON_KINESIS_VIDEO_STREAMS_FILE_NAME}"
        # Remove the downloaded file.
        rm -f "${AMAZON_KINESIS_VIDEO_STREAMS_TEMP_DOWNLOAD_PATH}"
        # Prepare the build directory.
        mkdir -p "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}"
        # Build the application.
        cmake -DBUILD_GSTREAMER_PLUGIN=TRUE -S "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_PATH}" -B "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}"
        # Build the application.
        make -C "${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}"
        # Add the path to the .bashrc file so that it can be used in the future
        echo -e "export GST_PLUGIN_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_PRODUCER_BUILD_PATH}:\$GST_PLUGIN_PATH\nexport LD_LIBRARY_PATH=${AMAZON_KINESIS_VIDEO_STREAMS_OPEN_SOURCE_LOCAL_LIB_PATH}:\$LD_LIBRARY_PATH" >>/root/.bashrc
        # Reload the .profile file.
        source /root/.bashrc
    fi
    # Check if the AWS CLI is installed.
    if [ ! -x "$(command -v aws)" ]; then
        # Install the AWS CLI
        curl "${AMAZON_CLI_LATEST_RELEASE}" -o "${AMAZON_CLI_TEMP_DOWNLOAD_PATH}"
        # Unzip the file
        unzip "${AMAZON_CLI_TEMP_DOWNLOAD_PATH}" -d /tmp/
        # Install the AWS CLI
        sudo ${AMAZON_CLI_INSTALL_SCRIPT_PATH}
        # Remove the downloaded file.
        rm -f "${AMAZON_CLI_TEMP_DOWNLOAD_PATH}"
        # Remove the downloaded file.
        rm -rf ${AMAZON_CLI_TEMP_INSTALL_PATH}
        # Login to aws. (https://us-east-1.console.aws.amazon.com/iamv2/home?region=us-east-1#/security_credentials)
        # aws configure set aws_access_key_id ${{ secrets.AWS_ACCESS_KEY }}
        # aws configure set aws_secret_access_key ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        # Run this command to manually feed data into Amazon Kinesis Video Streams
        # gst-launch-1.0 rtspsrc location=<rtsp_address> ! rtph264depay ! h264parse ! video/x-h264,stream-format=avc ! kvssink stream-name=<stream_id> access-key=<access_key> secret-key=<secret_key> aws-region=<aws_region>
    fi
}

# Build the application.
build-kensis-application

# Install Google Cloud
function install-google-cloud() {
    if { [ ! -x "$(command -v gcloud)" ] || [ ! -x "$(command -v vaictl)" ]; }; then
        # Install Google cloud sdk
        echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
        # Install the google cloud apt key
        curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
        apt-get update
        apt-get install google-cloud-cli -y
        # https://console.cloud.google.com/iam-admin/serviceaccounts
        # gcloud auth activate-service-account SERVICE_ACCOUNT@DOMAIN.COM --key-file=/path/key.json --project=PROJECT_ID
        # Install Google cloud vision ai
        curl -L "${GOOGLE_CLOUD_VISION_AI_LATEST_RELEASE}" -o "${GOOGLE_CLOUD_VISION_AI_TEMP_DOWNLOAD_PATH}"
        # Install the application from the downloaded file.
        apt-get install "${GOOGLE_CLOUD_VISION_AI_TEMP_DOWNLOAD_PATH}" -y
        # Remove the downloaded file.
        rm -f "${GOOGLE_CLOUD_VISION_AI_TEMP_DOWNLOAD_PATH}"
        # Run this command manually to feed data to GCP Vertex AI.
        # vaictl -p <project_id> -l <location_id> -c application-cluster-0 --service-endpoint visionai.googleapis.com send rtsp to streams <stream_id> --rtsp-uri <rtsp_address>
    fi
}

# Install Google Cloud
install-google-cloud


# Install the cloud connector.
function install-cps-connetor() {
    if [ ! -d "${CSP_CONNECTOR_PATH}" ]; then
        # Make the CSP connector directory
        mkdir -p "${CSP_CONNECTOR_PATH}"
        # Download the application
        curl -L "${CSP_CONNECTOR_LATEST_RELEASE}" -o "${CSP_CONNECTOR_APPLICATION}"
        # Download the config.
        curl -L "${CSP_CONNECTOR_CONFIG_URL}" -o "${CSP_CONNECTOR_CONFIG}"
        # Make the application executable
        chmod +x "${CSP_CONNECTOR_APPLICATION}"
        # Add the application to the path
        if [ ! -x "$(command -v csp-connector)" ]; then
            cp -s "${CSP_CONNECTOR_APPLICATION}" /usr/bin/csp-connector
        fi
        if [ ! -f "${CSP_CONNECTOR_SERVICE}" ]; then
            # This code creates the service file
            # The service file is stored in /etc/systemd/system/csp-connector.service
            echo "[Unit]
Description=Create the service for csp-connector.
Wants=network.target
[Service]
Type=simple
User=root
WorkingDirectory=${CSP_CONNECTOR_PATH}
ExecStart=csp-connector -config=\"${CSP_CONNECTOR_CONFIG}\" -log=\"${CSP_CONNECTOR_LOG_FILE}\" -twitch=true
Restart=always
[Install]
WantedBy=multi-user.target" >${CSP_CONNECTOR_SERVICE}
            # Reload the daemon
            systemctl daemon-reload
            if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
                systemctl enable --now csp-connector
                systemctl start csp-connector
            elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
                service csp-connector start
            fi
        fi
    fi
    # Update the latest version of the application
    if [ ! -f "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}" ]; then
        # Download the application
        curl -L "${CSP_CONNECTOR_LATEST_RELEASE}" -o "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}"
        # Make the application executable
        chmod +x "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}"
        # Check the sha512sum of the downloaded file
        CURRENT_FILE_SHA512SUM=$(sha512sum "${CSP_CONNECTOR_APPLICATION}" | awk '{print $1}')
        DOWNLOADED_FILE_SHA512SUM=$(sha512sum "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}" | awk '{print $1}')
        # Check if the sha512sum of the downloaded file is the same as the current file
        if [ "${CURRENT_FILE_SHA512SUM}" == "${DOWNLOADED_FILE_SHA512SUM}" ]; then
            # Remove the downloaded file.
            rm -f "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}"
        else
            # Stop the service
            if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
                systemctl stop csp-connector
            elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
                service csp-connector stop
            fi
            # Remove the application
            rm -f "${CSP_CONNECTOR_APPLICATION}"
            # Remove the application from the path
            if [ -x "$(command -v csp-connector)" ]; then
                rm -f ${CSP_CONNECTOR_BINARY_PATH}
            fi
            # Move the downloaded file to the application path
            mv "${CSP_CONNECTOR_TEMP_DOWNLOAD_PATH}" "${CSP_CONNECTOR_APPLICATION}"
            # Add the application to the path
            if [ ! -x "$(command -v csp-connector)" ]; then
                cp -s "${CSP_CONNECTOR_APPLICATION}" ${CSP_CONNECTOR_BINARY_PATH}
            fi
            # Start the service
            if [[ "${CURRENT_INIT_SYSTEM}" == *"systemd"* ]]; then
                systemctl start csp-connector
            elif [[ "${CURRENT_INIT_SYSTEM}" == *"init"* ]]; then
                service csp-connector start
            fi
        fi
    fi
}

# Install the cloud connector
install-cps-connetor
