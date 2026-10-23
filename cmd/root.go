package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jira-etl",
	Short: "Processes to extract data from Jira Rest API",
	Run: func(*cobra.Command, []string) {
		fmt.Println("Use subcommands")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
