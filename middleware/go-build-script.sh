# Build for all the OS
function build-golang-app() {
    APPLICATION="CSP-Connector"
    VERSION="v0.0.1"
    SOURCE_CODE="./main.go"
    BIN="bin/"
    if [ -n "$(ls ./*.go)" ]; then
        GOOS=linux GOARCH=386 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-386 ${SOURCE_CODE}
        GOOS=linux GOARCH=amd64 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-amd64 ${SOURCE_CODE}
        GOOS=linux GOARCH=arm go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm ${SOURCE_CODE}
        GOOS=linux GOARCH=arm64 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-arm64 ${SOURCE_CODE}
        GOOS=linux GOARCH=mips go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips ${SOURCE_CODE}
        GOOS=linux GOARCH=mips64 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64 ${SOURCE_CODE}
        GOOS=linux GOARCH=mips64le go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mips64le ${SOURCE_CODE}
        GOOS=linux GOARCH=mipsle go build -o ${BIN}${APPLICATION}-${VERSION}-linux-mipsle ${SOURCE_CODE}
        GOOS=linux GOARCH=ppc64 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64 ${SOURCE_CODE}
        GOOS=linux GOARCH=ppc64le go build -o ${BIN}${APPLICATION}-${VERSION}-linux-ppc64le ${SOURCE_CODE}
        GOOS=linux GOARCH=riscv64 go build -o ${BIN}${APPLICATION}-${VERSION}-linux-riscv64 ${SOURCE_CODE}
        GOOS=linux GOARCH=s390x go build -o ${BIN}${APPLICATION}-${VERSION}-linux-s390x ${SOURCE_CODE}
     fi
    else
        echo "Error: The \".go\" files could not be found."
        exit
    fi
}

build-golang-app
