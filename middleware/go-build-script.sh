# Define a function called build-golang-app
function build-golang-app() {
    # Navigate to the middleware directory
    cd middleware/
    # Set values for the application, version, source code and binaries
    APPLICATION="csp-connector"
    VERSION="v0.0.1"
    SOURCE_CODE="."
    BIN="binaries/"
    # Build the application for various OS and architectures using the go build command
    # Set the OS and architecture values for each build, and output the resulting binary to the BIN directory with the appropriate name
    GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-386 ${SOURCE_CODE}
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-amd64 ${SOURCE_CODE}
    GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm ${SOURCE_CODE}
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm64 ${SOURCE_CODE}
    GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips ${SOURCE_CODE}
    GOOS=linux GOARCH=mips64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64 ${SOURCE_CODE}
    GOOS=linux GOARCH=mips64le CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64le ${SOURCE_CODE}
    GOOS=linux GOARCH=mipsle CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mipsle ${SOURCE_CODE}
    GOOS=linux GOARCH=ppc64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64 ${SOURCE_CODE}
    GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64le ${SOURCE_CODE}
    GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-riscv64 ${SOURCE_CODE}
    GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-s390x ${SOURCE_CODE}
}
# Call the build-golang-app function to execute it
build-golang-app
