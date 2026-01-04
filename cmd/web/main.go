package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // New import
	"snippetbox.valak.gg/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	// Use the driver-compatible DSN format: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := flag.String("dsn", "root:RdARKKBtazsJtUvbQbfZPHgCeGPngSKa@tcp(hopper.proxy.rlwy.net:32223)/railway?parseTime=true", "MySQL data souce name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
