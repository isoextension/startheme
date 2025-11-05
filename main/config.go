package main

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Fuzzy       bool              `yaml:"fuzzy"`
	Sort        bool              `yaml:"sort"`
	StarshipDir string            `yaml:"starship_dir"`
	Editor      string            `yaml:"editor"`
	HintEnable  bool              `yaml:"hint_enable"`
	ThemeRepos  map[string]string `yaml:"theme_repos"`
}

func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		os_ := runtime.GOOS
		if os_ == "windows" {
			homeDir = os.Getenv("USERPROFILE")
		} else if os_ == "linux" {
			homeDir = os.Getenv("HOME")
		} else {
			return nil
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return &Config{
		Fuzzy:       false,
		Sort:        true,
		StarshipDir: filepath.Join(homeDir, ".config", "starship"),
		Editor:      editor,
		HintEnable:  true,
		ThemeRepos: map[string]string{
			"default": "git+https://startheme.github.io/default.git",
		},
	}
}