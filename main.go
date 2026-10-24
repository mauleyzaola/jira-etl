package main

import (
	_ "github.com/lib/pq"
	"github.com/mauleyzaola/jira-etl/cmd"
)

func main() {
	cmd.Execute()
}
