package main

import (
	"flag"
	"os"

	"github.com/mycok/shopit/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dsn string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "mongoDB DSN")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}