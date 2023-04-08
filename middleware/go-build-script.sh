# Build for all the OS
function build-golang-app() {
    cd middleware/
    APPLICATION="csp-connector"
    VERSION="v0.0.1"
    SOURCE_CODE="."
    BIN="binaries/"
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

build-golang-app
