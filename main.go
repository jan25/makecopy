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

func getRootDir(relpath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, relpath), nil
}

func walkFunc(changes []*Change) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fmt.Println(path)
		return nil
	}
}

func copyAndModify(filename string, changes []*Change) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	s := string(bytes)

	for _, change := range changes {
		old, new := change.Replace, change.With
		s = strings.ReplaceAll(s, old, new)
	}

	err = ioutil.WriteFile(filename, []byte(s), 0644)
	return err
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

	// validate config even if not errors

	fmt.Println(c.Message)

	// Make copy of target directory or target file
	// Clean copy if process was interupped

	// Ask for replacement tokens
	// Or use default tokens
	var changes []*Change
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
		changes = append(changes, change)
	}

	root, err := getRootDir(c.Path)
	if err != nil {
		log.Fatal(err)
	}
	if err := filepath.WalkDir(root, walkFunc(changes)); err != nil {
		log.Fatal(err)
	}

	// Modify files
	// Simple idea: Walk over files under a directory
	// Find and replace for one file at a time

	// Print new copy directory name
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
