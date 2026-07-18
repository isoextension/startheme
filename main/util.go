package main

import (
	"strings"

	"github.com/isoextension/btgo/ansi"
)

func hint(msg string, indent int) {
	indents := strings.Repeat(" ", indent)
	log.Plainf("%s%shint: %s%s", indents, ansi.BrightBlue.String(), msg, ansi.Reset.String())
}
