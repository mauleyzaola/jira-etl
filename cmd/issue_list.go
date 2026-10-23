package cmd

import (
	"fmt"
	"time"

	"github.com/mauleyzaola/jira-etl/envs"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
	"github.com/spf13/cobra"
)

var issueListProject string

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
	rootCmd.AddCommand(cmd)
}

func IssueList() error {
	config, err := envs.New()
	if err != nil {
		return err
	}
	from := time.Now().Add(-time.Hour * 24 * 30)
	to := time.Now()
	data, err := proxy.IssueList(config, issueListProject, from, to)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
