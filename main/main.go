package main

import (
	"runtime"

	"github.com/isoextension/btgo/taylog"
)

var log *taylog.Logger = taylog.New("startheme", nil)

func main() {
	debugf("Starting Startheme v%s with Go version %s\n", version, runtime.Version())
}
