go get -u github.com/c-darwin/dcoin-go-tmp
export CGO_ENABLED=1
export GOARCH=amd64 && go build -o dcoin64
export GOARCH=386 && go build -o dcoin32
zip dcoin_freebsd64.zip dcoin64
zip dcoin_freebsd32.zip dcoin32