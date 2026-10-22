package main

import (
	"fmt"
	"os"

	"github.com/mauleyzaola/jira-etl/envs"
	"github.com/mauleyzaola/jira-etl/internal/proxy"
)

func main() {
	if err := run(os.Args); err != nil {
		panic(err)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		fmt.Println("usage: jiralog key")
	}
	key := args[1]
	config, err := envs.New()
	if err != nil {
		return err
	}
	log, err := proxy.LoadJiraLog(config, key)
	if err != nil {
		return err
	}
	histories := proxy.ExtractStatuses(*log)
	for _, h := range histories {
		fmt.Println(h.Created, h.Author.DisplayName, *h.Items[0].FromString, *h.Items[0].ToString)
	}
	return nil
}
