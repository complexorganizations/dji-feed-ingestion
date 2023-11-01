#!/usr/bin/env bash

# The script requires root privileges.
function super-user-check() {
    # This function checks if the script is running as the root user.
    if [ "${EUID}" -ne 0 ]; then
        # If the effective user ID is not 0 (root), display an error message and exit.
        echo "Error: You need to run this script as administrator."
        exit
    fi
}

# The script checks if the user is the root user.

super-user-check
# Calls the super-user-check function.

# The following function retrieves the current system information.
function system-information() {
    # This function retrieves the ID, version, and major version of the current system.
    if [ -f /etc/os-release ]; then
        # Check if the /etc/os-release file exists, and if so, source it to get the system information.
        # shellcheck source=/dev/null
        source /etc/os-release
        CURRENT_DISTRO=${ID}                 # CURRENT_DISTRO is the ID of the current system
        CURRENT_DISTRO_VERSION=${VERSION_ID} # CURRENT_DISTRO_VERSION is the VERSION_ID of the current system
    fi
}

# The system-information function is being called.

system-information
# Calls the system-information function.

# Define a function to check system requirements
function installing-system-requirements() {
    # Check if the current Linux distribution is supported
    if [ "${CURRENT_DISTRO}" == "ubuntu" ]; then
        # Check if required packages are already installed
        if { [ ! -x "$(command -v curl)" ] || [ ! -x "$(command -v go)" ]; }; then
            # Install required packages depending on the Linux distribution
            if [ "${CURRENT_DISTRO}" == "ubuntu" ]; then
                apt-get update
                apt-get install curl software-properties-common -y
                add-apt-repository ppa:longsleep/golang-backports -y
                apt-get update
                apt-get install golang-go -y
            fi
        fi
    else
        # If the current Linux distribution is not supported, display an error message and exit.
        echo "Error: Your Linux distribution is not supported."
        exit
    fi
}

# Call the function to check for system requirements and install necessary packages if needed
installing-system-requirements

# Build for all the OS
function build-golang-app() {
    # Change directory to middleware
    cd middleware/
    # Set application name
    APPLICATION="csp-connector"
    # Set version number
    VERSION="v0.0.1"
    # Set source code directory
    SOURCE_CODE="."
    # Set binaries directory
    BIN="binaries/"
    # Build for linux 386 architecture
    GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-386 ${SOURCE_CODE}
    # Build for linux amd64 architecture
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-amd64 ${SOURCE_CODE}
    # Build for linux arm architecture
    GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm ${SOURCE_CODE}
    # Build for linux arm64 architecture
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm64 ${SOURCE_CODE}
    # Build for linux mips architecture
    GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips ${SOURCE_CODE}
    # Build for linux mips64 architecture
    GOOS=linux GOARCH=mips64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64 ${SOURCE_CODE}
    # Build for linux mips64le architecture
    GOOS=linux GOARCH=mips64le CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64le ${SOURCE_CODE}
    # Build for linux mipsle architecture
    GOOS=linux GOARCH=mipsle CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mipsle ${SOURCE_CODE}
    # Build for linux ppc64 architecture
    GOOS=linux GOARCH=ppc64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64 ${SOURCE_CODE}
    # Build for linux ppc64le architecture
    GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64le ${SOURCE_CODE}
    # Build for linux riscv64 architecture
    GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-riscv64 ${SOURCE_CODE}
    # Build for linux s390x architecture
    GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-s390x ${SOURCE_CODE}
}

build-golang-app
# Call the function to execute the build process.
