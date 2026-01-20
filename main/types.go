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

type Repo struct {
	Name string // specified by user
	URL  URL    // required
}

type Arguments struct {
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
