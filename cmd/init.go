package cmd

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

// Initial configuration
// Keep this in sync with ../.makecopy.yml file
var initialConfig string = `# Target directory to copy
path: "./example/example-app"

# Not supported yet: git url
git: "github.com/path/to/repo"

# Optional message to user
# Default is "Copying {path}"
message: "Copying example-app"

# Changes to apply to copied files
changes:
  # Helper message about change
  "Project prefix":
    # Replace this token
    replace: "example"
    # With optional default value
    default: "myexample"

  "App name":
    replace: "app"
    default: "myapp"

# Not supported yet: Change file and directory names too, not just the content
changefilenames: true`

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration.",
	Long:  `Setup a sample configuration for copying code files.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ioutil.WriteFile(".makecopy.yml", []byte(initialConfig), 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}
