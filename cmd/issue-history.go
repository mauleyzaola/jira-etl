package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mauleyzaola/jira-etl/envs"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
	"github.com/spf13/cobra"
)

var issueHistoryCmd = &cobra.Command{
	Use:   "issue-history",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing issue key")
		}
		config, err := envs.New()
		if err != nil {
			return err
		}
		return IssueHistory(config, args[0])
	},
}

func init() {
	cmd := issueHistoryCmd
	rootCmd.AddCommand(cmd)
}

func IssueHistory(config envs.Config, key string) error {
	out, err := proxy.IssueHistory(config, key)
	if err != nil {
		return err
	}
	data, err := json.Marshal(out)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
