package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/isoextension/btgo/logger"
)

var l = logger.New("startheme", nil)

func main() {
	var homeDir, err = os.UserHomeDir()
	var starshipDir = filepath.Join(homeDir, ".config", "starship")
	var starshipConfig = filepath.Join(homeDir, ".config", "starship.toml")

	if err != nil {
		l.BasicErrorf("error getting $HOME. %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "list":
		themes, err := getThemes(starshipDir)
		if err != nil {
			l.BasicErrorf("error listing themes, %v", err)
			os.Exit(1)
		}
		for _, theme := range themes {
			fmt.Println("  " + theme.name)
		}

	case "change", "use":
		if len(os.Args) < 3 {
			l.BasicErrorf("arguments missing")
			os.Exit(1)
		}
		err := changeTheme(starshipDir, starshipConfig, os.Args[2])
		if err != nil {
			l.BasicErrorf("error changing theme: %v", err)
			os.Exit(1)
		}

	case "get":
		theme, err := getCurrentTheme(starshipConfig)
		if err != nil {
			l.BasicErrorf("error getting current theme: %v", err)
			os.Exit(1)
		}
		fmt.Println(theme.name)

	case "help", "-h", "--help":
		showHelp()

	case "edit", "open":
		theme2Open, err := getCurrentTheme(starshipConfig)
		if err != nil {
			l.BasicErrorf("%v", err)
			os.Exit(1)
		}
		err = openTheme(*theme2Open, starshipConfig)
		if err != nil {
			l.BasicErrorf("error opening theme: %v", err)
			os.Exit(1)
		}

	default:
		l.BasicErrorf("no such command %s", os.Args[1])
		showHelp()
		os.Exit(1)
	}
}