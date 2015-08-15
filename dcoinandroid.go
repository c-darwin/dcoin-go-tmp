// +build android

package main

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcoin"
	"golang.org/x/mobile/app"
)

/*
#include <stdio.h>
#include <stdlib.h>

char* JGetTmpDir2() {
	return getenv("TMPDIR");
}
*/
import "C"

func main() {
	app.Main(func(a app.App) {
		dcoin.Start(C.GoString(C.JGetTmpDir2()))
	})
}