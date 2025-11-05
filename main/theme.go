package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

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

func getThemes(starshipDir string) ([]*theme, error) {
	entries, err := os.ReadDir(starshipDir)
	if err != nil {
		return nil, err
	}
	var themes []*theme
	for _, entry := range entries {
		var themePath = filepath.Join(starshipDir, entry.Name())
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".toml" {
			themeName := entry.Name()[:len(entry.Name())-5]
			themes = append(themes, newtheme(themePath, themeName))
		} else if entry.IsDir() {
			themeName := entry.Name() + "/"
			themes = append(themes, newtheme(themePath, themeName))
		}
	}
	return themes, nil
}

func lookupTheme(substr string, starshipDir string, sorted bool) ([]*theme, error) {
	themes, err := getThemes(starshipDir)
	if err != nil {
		return nil, err
	}
	var themes2Return []*theme

	for _, theme := range themes {
		if strings.Contains(strings.ToLower(theme.name), strings.ToLower(substr)) {
			themes2Return = append(themes2Return, theme)
		}
	}

	if sorted {
		sort.SliceStable(themes2Return, func(i, j int) bool {
			return themes2Return[i].name < themes2Return[j].name
		})
	}

	return themes2Return, nil
}

func changeTheme(starshipDir, starshipConfig, themeName string) error {
	themeFile := filepath.Join(starshipDir, themeName+".toml")
	info, err := os.Stat(themeFile)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return syscall.EISDIR
	}
	if _, err := os.Lstat(starshipConfig); err == nil {
		if err := os.Remove(starshipConfig); err != nil {
			return err
		}
	}
	if err := os.Symlink(themeFile, starshipConfig); err != nil {
		return err
	}
	return nil
}

func getCurrentTheme(starshipConfig string) (*theme, error) {
	info, err := os.Lstat(starshipConfig)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(starshipConfig)
		if err != nil {
			return nil, err
		}
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

func openTheme(target theme, starshipConfig string) error {
	var editor string
	if target.path == "" {
		currentTheme, err := getCurrentTheme(starshipConfig)
		if err != nil {
			return err
		}
		target = *currentTheme
	}
	editor = os.Getenv("EDITOR")
	if editor == "" {
		l.Warning("$EDITOR is empty, defaulting to vim")
		l.Info("hint: try doing this:\n  EDITOR=myeditor startheme edit")
		editor = "vim"
	}
	var editorproc = exec.Command(editor, target.path)
	editorproc.Stdin = os.Stdin
	editorproc.Stdout = os.Stdout
	editorproc.Stderr = os.Stderr
	return editorproc.Run()
}