package internal

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileOp wraps a file operation type
type FileOp interface {
	Operate(path string, isDir bool) error
}

// Copier operated to copy files recursively
type Copier struct {
	Source string
	Dest   string
}

// Operate copies files
func (c *Copier) Operate(sourcePath string, isDir bool) error {
	rel, err := filepath.Rel(c.Source, sourcePath)
	if err != nil {
		return err
	}

	destpath := filepath.Join(c.Dest, rel)

	if isDir {
		return os.Mkdir(destpath, 0755)
	}

	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destpath, data, 0655)
}

// Changer changes files after copying
type Changer struct {
	Changes []*Change
}

// Operate changes files based on config
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

// WalkFunc wraps fs.WalkDirFunc with a specific FileOp.
func WalkFunc(op FileOp) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		return op.Operate(path, d.IsDir())
	}
}

// GetSourceDir computes absolute path to source directory.
func GetSourceDir(relpath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, relpath), nil
}

// GetDestDir computes absolute path to destination directory.
func GetDestDir(relpath string) (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	olddir := filepath.Base(relpath)
	newdir := fmt.Sprintf("%s-%s", olddir, randomName())
	return filepath.Join(d, newdir), nil
}
