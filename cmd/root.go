package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jan25/makecopy/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "makecopy",
	Aliases: []string{"mc"},
	Short:   "Make a copy of code",
	Long: `Make copy of example code and replace pieces of it
as you prefer with predefined configuration.`,

	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		// Say that we're copying files
		fmt.Println(config.Message)

		// Copy files from source directory
		from, err := internal.GetSourceDir(config.Path)
		if err != nil {
			log.Fatal(err)
		}
		to, err := internal.GetDestDir(config.Path)
		if err != nil {
			log.Fatal(err)
		}

		os.Mkdir(to, 0755)
		copier := internal.Copier{
			Source: from,
			Dest:   to,
		}
		copierFunc := internal.WalkFunc(&copier)
		if err := filepath.WalkDir(from, copierFunc); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Cloning done")

		// Change copied files
		fmt.Println("Changing copied files")
		changer := internal.Changer{}

		for question, change := range config.Changes {
			if change.With == "" {
				answer, err := internal.Prompt(question, change.Default)
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

		changerFunc := internal.WalkFunc(&changer)
		if err := filepath.WalkDir(to, changerFunc); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Changing done")

		// Print new copy directory name
		fmt.Printf("Success making copy of %s at %s \n", config.Path, filepath.Base(to))
	},
}

// Execute executes root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
