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
