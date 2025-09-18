package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	colorRed   = "\033[1;31m"
	colorGreen = "\033[1;32m"
	colorBlue  = "\033[1;34m"
	colorReset = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%s✖ Error getting home directory: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	starshipDir := filepath.Join(homeDir, ".config", "starship")
	starshipConfig := filepath.Join(homeDir, ".config", "starship.toml")

	switch os.Args[1] {
	case "list":
		listThemes(starshipDir)
	case "change":
		if len(os.Args) < 3 {
			fmt.Printf("%s✖ missing arguments%s\n", colorRed, colorReset)
			os.Exit(1)
		}
		changeTheme(starshipDir, starshipConfig, os.Args[2])
	case "get":
		getCurrentTheme(starshipConfig)
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("%s✖ Unknown subcommand: %s%s\n", colorRed, os.Args[1], colorReset)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	helpText := `startheme
    a Go program for hotswapping starship themes
    
    usage: startheme <subcommand>

    available subcommands:
        change   Change the current theme
        list     List themes
        get      Get current theme
        help     Show this help message

    WARNING: this command isn't going to back up your current starship theme. so please do so
`
	fmt.Print(helpText)
}

func listThemes(starshipDir string) {
	entries, err := os.ReadDir(starshipDir)
	if err != nil {
		fmt.Printf("%s✖ Error reading starship directory: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("%sAvailable themes:%s\n", colorBlue, colorReset)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".toml" {
			// Remove .toml extension for cleaner output
			themeName := entry.Name()[:len(entry.Name())-5]
			fmt.Printf("  %s\n", themeName)
		}
	}
}

func changeTheme(starshipDir, starshipConfig, themeName string) {
	themeFile := filepath.Join(starshipDir, themeName+".toml")

	// Check if theme file exists and is not a directory
	info, err := os.Stat(themeFile)
	if os.IsNotExist(err) {
		fmt.Printf("%s✖ %s: No such file or directory%s\n", colorRed, themeName, colorReset)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("%s✖ Error checking theme file: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Printf("%s✖ %s is a directory%s\n", colorRed, themeName, colorReset)
		os.Exit(1)
	}

	// Remove existing symlink/file if it exists
	if _, err := os.Lstat(starshipConfig); err == nil {
		if err := os.Remove(starshipConfig); err != nil {
			fmt.Printf("%s✖ Error removing existing config: %v%s\n", colorRed, err, colorReset)
			os.Exit(1)
		}
	}

	// Create new symlink
	if err := os.Symlink(themeFile, starshipConfig); err != nil {
		fmt.Printf("%s✖ Error creating symlink: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("%s✓ Successfully changed theme to: %s%s\n", colorGreen, themeName, colorReset)
}

func getCurrentTheme(starshipConfig string) {
	// Check if starship.toml exists and is a symlink
	info, err := os.Lstat(starshipConfig)
	if os.IsNotExist(err) {
		fmt.Printf("%s✖ No starship config found%s\n", colorRed, colorReset)
		return
	}
	if err != nil {
		fmt.Printf("%s✖ Error checking config: %v%s\n", colorRed, err, colorReset)
		return
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(starshipConfig)
		if err != nil {
			fmt.Printf("%s✖ Error reading symlink: %v%s\n", colorRed, err, colorReset)
			return
		}

		// Extract theme name from path
		themeName := filepath.Base(target)
		if filepath.Ext(themeName) == ".toml" {
			themeName = themeName[:len(themeName)-5]
		}

		fmt.Printf("%sCurrent theme: %s%s\n", colorBlue, themeName, colorReset)
	} else {
		fmt.Printf("%s✖ starship.toml is not a symlink (not managed by startheme)%s\n", colorRed, colorReset)
	}
}
