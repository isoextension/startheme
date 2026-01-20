package main

import (
	"os"

	"github.com/go-git/go-git"
)

var repos []Repo

func clone(repo *Repo, dir string) (*git.Repository, error) {
	err := os.MkdirAll(dir, 0o700)
	if err != nil {
		return nil, err
	}
	return git.PlainClone(dir, false, &git.CloneOptions{URL: string(repo.URL)})
}
