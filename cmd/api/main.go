package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/mycok/shopit/internal/db/mongo"
	"github.com/mycok/shopit/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port     int
	env      string
	dbClient struct {
		dsn string
	}
}

type application struct {
	config       config
	logger       *jsonlog.Logger
	wg           sync.WaitGroup
	repositories *mongo.Repositories
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.dbClient.dsn, "db-dsn", "", "MongoDB_DSN")
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
		repositories: mongo.NewRepositories(db.DB),
	}

	err = app.serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
