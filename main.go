package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	config := getConfig()

	// Say that we're copying files
	fmt.Println(config.Message)

	// Copy files from source directory
	from, err := getSourceDir(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	to, err := getDestDir(config.Path)
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(to, 0755)
	copier := Copier{
		Source: from,
		Dest:   to,
	}
	copierFunc := walkFunc(&copier)
	if err := filepath.WalkDir(from, copierFunc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cloning done")

	// Change copied files
	fmt.Println("Changing copied files")
	changer := Changer{}

	for question, change := range config.Changes {
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
	if err := filepath.WalkDir(to, changerFunc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Changing done")

	// Print new copy directory name
	fmt.Printf("Success making copy of %s at %s \n", config.Path, filepath.Base(to))
}
