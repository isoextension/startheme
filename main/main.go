package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/isoextension/btgo/ansi"
	"github.com/isoextension/btgo/logger"
)

var homeDir, hmdirerr = os.UserHomeDir()
var starshipDir = filepath.Join(homeDir, ".config", "starship")
var starshipConfig = filepath.Join(homeDir, ".config", "starship.toml")
var l *logger.Logger = logger.New("startheme", nil)

func hint(hint string) {
    if os.Getenv("STARTHEME_HINT_ENABLE") == "true" ||
        os.Getenv("STARTHEME_HINT_ENABLE") == "1" {
        l.Fplainf(os.Stderr, "%s%s%s\n", ansi.BrightBlack.String(), hint, ansi.Reset.String())
    }
}

type theme struct {
    path string
    name string
}

func newtheme(path string, name string) *theme {
    return &theme{
        path: path,
        name: name,
    }
}

func main() {
    if hmdirerr != nil {
        panic(hmdirerr)
    }
    if len(os.Args) < 2 {
        showHelp()
        return
    }
    switch os.Args[1] {
    case "list":
        listThemes(starshipDir)
    case "change", "use":
        if len(os.Args) < 3 {
            l.BasicErrorf("arguments missing")
            os.Exit(1)
        }
        changeTheme(starshipDir, starshipConfig, os.Args[2])
    case "get":
        getCurrentTheme(starshipConfig)
    case "help", "-h", "--help":
        showHelp()
	case "edit":
		theme2Open, err := getCurrentTheme(starshipConfig)
		l.BasicErrorf("%v", err)
		openTheme(*theme2Open)
    default:
        l.BasicErrorf("no such command %s", os.Args[2])
        showHelp()
        os.Exit(1)
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
        edit               Edit current theme
WARNING: this command isn't going to back up your current starship theme. so please do so
`
    fmt.Print(helpText)
}

func getThemes(starshipDir string) ([]string, error) {
    entries, err := os.ReadDir(starshipDir)
    if err != nil {
        return nil, err
    }
    var themes []string
    for _, entry := range entries {
        if !entry.IsDir() && filepath.Ext(entry.Name()) == ".toml" {
            // Remove .toml extension for cleaner output
            themeName := entry.Name()[:len(entry.Name())-5]
            themes = append(themes, themeName)
        }
    }
    return themes, nil
}

func listThemes(starshipDir string) {
    themes, err := getThemes(starshipDir)
    if err != nil {
        l.Fplainf(os.Stderr, "%s✖ Error reading starship directory: %v%s\n", ansi.Red.String(), err, ansi.Reset.String())
        os.Exit(1)
    }
    for _, theme := range themes {
        fmt.Printf("  %s\n", theme)
    }
}

func changeTheme(starshipDir, starshipConfig, themeName string) error {
    themeFile := filepath.Join(starshipDir, themeName+".toml")
    // Check if theme file exists and is not a directory
    info, err := os.Stat(themeFile)
    if err != nil {
        return err
    }
    if info.IsDir() {
        return syscall.EISDIR
    }
    // Remove existing symlink/file if it exists
    if _, err := os.Lstat(starshipConfig); err == nil {
        if err := os.Remove(starshipConfig); err != nil {
            return err
        }
    }
    // Create new symlink
    if err := os.Symlink(themeFile, starshipConfig); err != nil {
        return err
    }
    return nil
}

func getCurrentTheme(starshipConfig string) (*theme, error) {
    // Check if starship.toml exists and is a symlink
    info, err := os.Lstat(starshipConfig)
    if err != nil {
        return nil, err
    }
    if info.Mode()&os.ModeSymlink != 0 {
        target, err := os.Readlink(starshipConfig)
        if err != nil {
            return nil, err
        }
        // Extract theme name from path
        var theme string = filepath.Base(target)
        var themeName string
        var ext string = filepath.Ext(theme)
        var extlen int = len(ext)
        if extlen <= 0 {
            themeName = theme
        } else {
            themeName = theme[:len(theme)-extlen]
        }
        return newtheme(target, themeName), nil
    } else {
        return nil, errors.New("STARTHEME_NOT_MANAGED")
    }
}

func openTheme(target theme) error {
    var editor string
    if target.path == "" {
        currentTheme, err := getCurrentTheme(starshipConfig)
        if err != nil {
            // handle error
            return err
        }
        target = *currentTheme
    }
    editor = os.Getenv("EDITOR")
    if editor == "" {
        l.Warning("$EDITOR is empty, defaulting to /bin/env vim")
        l.Info("hint: try doing this:\n  EDITOR=myeditor startheme edit")
        editor = "/bin/env vim"
    }
    var editorproc *exec.Cmd = exec.Command(editor, target.path)
    var err = editorproc.Run()
    if err != nil {
        return err
    }
    return nil
}