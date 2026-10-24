package cmd

import (
	"fmt"
	"os"

	"github.com/mauleyzaola/jira-etl/internal/envs"
	"github.com/mauleyzaola/jira-etl/internal/formatters/jsons"
	"github.com/mauleyzaola/jira-etl/internal/ports"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
	"github.com/spf13/cobra"
)

var issueHistoryCmd = &cobra.Command{
	Use:   "issue-history",
	Short: "Returns a list of issue history (log) for a given issue",
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
	cmd.Flags().StringVar(&format, "format", "", "json")
	_ = cmd.MarkFlagRequired("format")
	rootCmd.AddCommand(cmd)
}

func IssueHistory(config envs.Config, key string) error {
	allowedFormatters := map[string]ports.Formatter{
		"json": jsons.NewGenericFormatter(),
	}
	formatter, ok := allowedFormatters[format]
	if !ok {
		return fmt.Errorf("formatter not supported: %s", format)
	}
	return proxy.IssueHistory(os.Stdout, formatter, config, key)
}
