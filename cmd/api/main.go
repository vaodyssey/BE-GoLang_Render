package main

import (
	db "TechStore/db/sqlc"
	"TechStore/internal/pkg/jsonlog"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	port int

	db struct {
		dsn              string
		maxConnsLifeTime time.Duration
		maxOpenConns     int
		maxIdleConns     int
	}
}

type application struct {
	config  config
	logger  *jsonlog.Logger
	queries *db.Queries
	db      *sql.DB
	wg      sync.WaitGroup
}

const (
	Env     = ".env"
	AppPort = "APP_PORT"

	DbDsn = "DB_DSN"
)

func main() {
	err := godotenv.Load(Env)
	logr := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	if err != nil {
		logr.PrintFatal(err, nil)
	}

	var cfg config

	// Load port configurations
	port, err := strconv.Atoi(os.Getenv(AppPort))
	if err != nil {
		logr.PrintFatal(err, nil)
	}
	cfg.port = port

	// Load sqlcDb configurations
	cfg.db.dsn = os.Getenv(DbDsn)
	cfg.db.maxConnsLifeTime = time.Minute * 3
	cfg.db.maxOpenConns = 20
	cfg.db.maxIdleConns = 20

	// Initialize sqlcDb
	dbConn, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		logr.PrintFatal(err, nil)
	}

	dbConn.SetConnMaxLifetime(cfg.db.maxConnsLifeTime)
	dbConn.SetMaxOpenConns(cfg.db.maxOpenConns)
	dbConn.SetMaxIdleConns(cfg.db.maxIdleConns)

	app := &application{
		config:  cfg,
		logger:  logr,
		db:      dbConn,
		queries: db.New(dbConn),
		wg:      sync.WaitGroup{},
	}

	if err := app.serve(); err != nil {
		logr.PrintFatal(err, nil)
	}

}
