package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileOp interface {
	Operate(path string, isDir bool) error
}

type Copier struct {
	Source string
	Dest   string
}

func (c *Copier) Operate(path string, isDir bool) error {
	rel, err := filepath.Rel(c.Source, path)
	if err != nil {
		return err
	}

	destpath := filepath.Join(c.Dest, rel)

	if isDir {
		return os.Mkdir(destpath, 0755)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destpath, data, 0655)
}

type Changer struct {
	Changes []*Change
}

func (c *Changer) Operate(path string, isDir bool) error {
	if isDir {
		return nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	s := string(bytes)

	for _, change := range c.Changes {
		old, new := change.Replace, change.With
		s = strings.ReplaceAll(s, old, new)
	}
	return ioutil.WriteFile(path, []byte(s), 0644)
}
