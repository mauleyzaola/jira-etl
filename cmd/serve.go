package cmd

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauleyzaola/jira-etl/internal/domains"
	"github.com/mauleyzaola/jira-etl/internal/envs"
	"github.com/mauleyzaola/jira-etl/internal/helpers"
	"github.com/mauleyzaola/jira-etl/internal/stores/postgres"
	"github.com/spf13/cobra"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts HTTP server to listen for jira webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Serve(); err != nil {
			panic(err)
		}
	},
}

func init() {
	cmd := serveCmd
	cmd.Flags().StringVar(&port, "port", "8080", "http server listen port")
	rootCmd.AddCommand(cmd)
}

// jira webhook secret
// TXl89l7PBwa4Ltzur1ZU

// findings
// webhook timestamp is expressed in milliseconds

func Serve() error {
	ctx := context.Background()
	config, err := envs.New()
	if err != nil {
		return err
	}
	db, err := helpers.PoolConnectFromEnvVars(ctx, config.PostgresConfig)
	if err != nil {
		return err
	}
	defer func() { db.Close() }()
	store := postgres.NewStore()
	router := chi.NewRouter()
	router.Post("/jira/webhooks", func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming jira webhook request")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wh, err := domains.NewJiraWebhook(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = store.InsertJiraWebhook(ctx, db, *wh); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Println("listening on :" + port)
	return http.ListenAndServe(":"+port, router)
}
