#! /bin/bash -e

go get -u github.com/c-darwin/dcoin-go-tmp
GOARCH=386  CGO_ENABLED=1  go build -o dcoin/usr/share/dcoin/dcoin
dpkg-deb --build dcoin
go get -u github.com/c-darwin/dcoin-go-tmp
GOARCH=amd64  CGO_ENABLED=1  go build -o dcoin64/usr/share/dcoin/dcoin
dpkg-deb --build dcoin64