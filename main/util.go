package main

import (
	"fmt"
	"os"

	"github.com/isoextension/btgo/ansi"
)

func hint(hint string) {
	if os.Getenv("STARTHEME_HINT_ENABLE") == "true" ||
		os.Getenv("STARTHEME_HINT_ENABLE") == "1" {
		l.Fplainf(os.Stderr, "%s%s%s\n", ansi.BrightBlack.String(), hint, ansi.Reset.String())
	}
}

func showHelp() {
	helpText := `
startheme - a Go program for hotswapping starship themes
usage: startheme <subcommand> <args>
    available subcommands:
        change, use        Change the current theme
        list               List themes
        get                Get current theme
        help, -h, --help   Show this help message
        edit               Edit theme based on argument passed
                           if no theme was passed then edits current theme

WARNING: this command isn't going to back up your current starship theme. so please do so
`
	fmt.Print(helpText)
}