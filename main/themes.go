package main

import "os"

func switchTheme(theme *Theme, dir string) error {
	_, err := os.Stat(theme.Path())
	if err != nil {
		return err
	}
	return os.Symlink(theme.Path(), dir)
}
