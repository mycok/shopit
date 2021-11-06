package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/mycok/shopit/internal/db/mongo"
	"github.com/mycok/shopit/internal/jsonlog"
	"github.com/mycok/shopit/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port     int
	env      string
	dbClient struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config       config
	logger       *jsonlog.Logger
	wg           sync.WaitGroup
	mailer       mailer.Mailer
	repositories *mongo.Repositories
}

func main() {
	var cfg config

	// Server connection config.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.dbClient.dsn, "db-dsn", "", "MongoDB_DSN")

	// Mail SMTP config.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "b5054223ea31d2", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "967bce52ce4e76", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "<no-reply@shopit.co.ug>", "SMTP sender")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Open a mongo server connection.
	client, err := mongo.OpenConnection(cfg.dbClient.dsn)
	if err != nil {
		logger.LogFatal(err, nil)
	}

	// Create a database and register new collections.
	db := mongo.New(client, "shopit")

	err = db.RegisterNewCollections()
	if err != nil {
		logger.LogFatal(err, nil)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.LogFatal(err, nil)
		}
	}()

	app := &application{
		config:       cfg,
		logger:       logger,
		mailer:       mailer.New(cfg.smtp.port, cfg.smtp.host, cfg.smtp.sender, cfg.smtp.username, cfg.smtp.password),
		repositories: mongo.NewRepositories(db.DB),
	}

	err = app.serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
