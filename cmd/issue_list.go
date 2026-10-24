package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mauleyzaola/jira-etl/internal/envs"
	"github.com/mauleyzaola/jira-etl/internal/formatters/jsons"
	"github.com/mauleyzaola/jira-etl/internal/ports"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	issueListProject string
	fromStr, toStr   string
	format           string
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
	cmd.Flags().StringVar(&format, "format", "", "json")
	_ = cmd.MarkFlagRequired("format")
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
	allowedFormatters := map[string]ports.Formatter{
		"json": jsons.NewGenericFormatter(),
	}
	formatter, ok := allowedFormatters[format]
	if !ok {
		return fmt.Errorf("formatter not supported: %s", format)
	}
	w := os.Stdout
	return proxy.IssueList(w, formatter, config, issueListProject, from, to)
}
