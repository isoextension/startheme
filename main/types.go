package main

import (
	"fmt"
)

// / ALIASES ///
type Code []byte
type URL string

// / STRUCTS ///
type Theme struct {
	name string
	path string // required
	Code        // required
}

type SwitchCmd struct {
	Name      string `arg:"positional"`
	RepoToGet string `arg:"-r,--repo,env:STARTHEME_REPO"`
}

type GetCurrent struct {
	ShowCode bool `arg:"-c,--code,--show-code"`
	ShowPath bool `arg:"-p,--path,--show-path"`
}

type ListCmd struct{}

type Arguments struct {
	Switch *SwitchCmd  `arg:"subcommand:use"`
	Get    *GetCurrent `arg:"subcommand:get"`
	List   *ListCmd    `arg:"subcommand:list"`
}

// / METHODS ///
func (c Code) String() string        { return string(c) }
func (t Theme) Name() string         { return t.name }
func (t Theme) Path() string         { return t.path }
func (t *Theme) SetName(name string) { t.name = name }
func (t *Theme) SetPath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	t.path = path
	return nil
}

// / CONSTRUCTORS ///
func NewTheme(name string, path string, code Code) (*Theme, error) {
	if path == "" {
		return nil, fmt.Errorf("no path was specified")
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no code was specified")
	}
	return &Theme{name: name, path: path, Code: code}, nil
}
