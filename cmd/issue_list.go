package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mauleyzaola/jira-etl/envs"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	issueListProject string
	fromStr, toStr   string
)

var issueListCmd = &cobra.Command{
	Use:   "issue-list",
	Short: "Returns a list of jira issues which have been updated within a date rage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return IssueList()
	},
}

func init() {
	cmd := issueListCmd
	cmd.Flags().StringVar(&issueListProject, "project", "", "Jira project key")
	_ = cmd.MarkFlagRequired("project")
	cmd.Flags().StringVar(&fromStr, "from", "", "Start date format 2006-01-02")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&toStr, "to", "", "End date format 2006-01-02")
	_ = cmd.MarkFlagRequired("to")
	rootCmd.AddCommand(cmd)
}

func IssueList() error {
	const dateFormat = "2006-01-02"
	config, err := envs.New()
	if err != nil {
		return err
	}
	from, err := time.Parse(dateFormat, fromStr)
	if err != nil {
		return err
	}
	to, err := time.Parse(dateFormat, toStr)
	if err != nil {
		return err
	}
	issues, err := proxy.IssueList(config, issueListProject, from, to)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
