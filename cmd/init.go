package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration.",
	Long:  `Setup a sample configuration for copying code files.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Not implemented")
	},
}
