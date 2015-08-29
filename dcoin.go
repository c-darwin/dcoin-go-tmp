// +build !android,!ios

package main

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcoin"
)

func main() {
	dcoin.Start("")
}