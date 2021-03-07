package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Change struct {
	Replace string `yaml:"replace"`
	With    string `yaml:"with"`
	Default string `yaml:"default"`
}

type Config struct {
	Message         string             `yaml:"message"`
	Path            string             `yaml:"path"`
	Changes         map[string]*Change `yaml:"changes"`
	ModifyFilenames bool               `yaml:"modifyfilenames"`
}

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

func walkFunc(op FileOp) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		return op.Operate(path, d.IsDir())
	}
}

func getSourceDir(relpath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, relpath), nil
}

func getDestDir(relpath string) (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	olddir := filepath.Dir(d)
	newdir := fmt.Sprintf("%s-%s", olddir, randomName())
	return filepath.Join(d, newdir), nil
}

func main() {
	bytes, err := ioutil.ReadFile("./.makecopy.yml")
	if err != nil {
		log.Fatal(err)
	}

	c := Config{}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		log.Fatal(err)
	}

	// Validate config

	// Say that we're copying files
	fmt.Println(c.Message)

	// Copy files from source directory
	from, err := getSourceDir(c.Path)
	if err != nil {
		log.Fatal(err)
	}
	to, err := getDestDir(c.Path)
	if err != nil {
		log.Fatal(err)
	}

	copier := Copier{
		Source: from,
		Dest:   to,
	}
	copierFunc := walkFunc(&copier)
	if err := filepath.WalkDir(from, copierFunc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cloning done")

	// Change copied over files
	fmt.Println("Changing copied files")
	changer := Changer{}

	for question, change := range c.Changes {
		if change.With == "" {
			answer, err := prompt(question, change.Default)
			if err != nil {
				log.Fatal(err)
			}
			change.With = answer
		}
		if change.With == "" {
			change.With = change.Default
		}
		changer.Changes = append(changer.Changes, change)
	}

	changerFunc := walkFunc(&changer)
	if err := filepath.WalkDir(from, changerFunc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Changing done")

	// Print new copy directory name
	fmt.Printf("Success making copy of %s at %s \n", c.Path, filepath.Dir(to))
}

func prompt(q string, suggest string) (string, error) {
	p := ""
	if suggest != "" {
		p = fmt.Sprintf("%s (%s): ", q, suggest)
	} else {
		p = fmt.Sprintf("%s: ", q)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	result, err := reader.ReadString('\n')
	return result, err
}
