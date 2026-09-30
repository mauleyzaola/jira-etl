package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loadTicketCmd = &cobra.Command{
	Use:   "load-ticket",
	Short: "Reads one jira ticket from Rest API and returns the data as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("load-ticket called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadTicketCmd)
}
