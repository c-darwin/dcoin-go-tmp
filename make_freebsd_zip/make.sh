go get -u github.com/c-darwin/dcoin-go-tmp
setenv GOARCH amd64 && setenv CGO_ENABLED 1 && go build -o dcoin64
setenv GOARCH 386 && setenv CGO_ENABLED 1 && go build -o dcoin32
zip dcoin_freebsd64.zip dcoin64
zip dcoin_freebsd32.zip dcoin32