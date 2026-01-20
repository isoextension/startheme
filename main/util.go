package main

import (
	"bytes"
	"fmt"
)

var debugbuf bytes.Buffer

func debugf(format string, args ...interface{}) {
	log.Fplainf(&debugbuf, format, args...)
}

func debugln(args ...interface{}) {
	fmt.Fprintln(&debugbuf, args...)
}
