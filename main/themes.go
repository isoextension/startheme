package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotManaged = fmt.Errorf("TOML file is not managed by startheme. Please run the \"use\" command first")

func switchTheme(theme *Theme, dir string) error {
	err := os.Remove(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	_, err = os.Stat(theme.Path())
	if err != nil {
		return err
	}
	return os.Symlink(theme.Path(), dir)
}

func current(file string) (*Theme, error) {
	info, err := os.Lstat(file)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return nil, ErrNotManaged
	}

	path, err := os.Readlink(file)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))

	return &Theme{
		name: name,
		path: path,
		Code: data,
	}, nil
}

func lookup(name string, dir string) (*Theme, error) {
	filename := fmt.Sprintf("%s.toml", name)
	path := filepath.Join(dir, filename)
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Theme{
		name: name,
		path: path,
		Code: data,
	}, nil
}
