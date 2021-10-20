package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/mycok/shopit/internal/jsonlog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func openDB(logger *jsonlog.Logger, cfg config) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.dbClient.dsn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.LogFatal(err, nil)
		}
	}()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.dbClient.dsn, "db-dsn", "", "MongoDB_DSN")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	_, err := openDB(app.logger, cfg)
	if err != nil {
		logger.LogFatal(err, nil)
	}

	err = app.serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
