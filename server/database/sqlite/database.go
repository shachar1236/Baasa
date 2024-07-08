package sqlite

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shachar1236/Baasa/database/sqlite/objects"
)

//go:embed schema.sql
var ddl string

type SqlTypes interface {
	int64 | float64 | string | []byte
}

type SqliteDB struct {
	db              *sqlx.DB
	objects_queries *objects.Queries

	logger    *slog.Logger
}

func New(ctx context.Context) (SqliteDB, error) {
	logFile, err := os.OpenFile("logs/sqlite_db.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	var mw io.Writer
	if err != nil {
		mw = os.Stdout
	} else {
		mw = io.MultiWriter(os.Stdout, logFile)
	}

	var this SqliteDB
	this.logger = slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))

	mydb, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		return this, err
	}

	// create tables
	if _, err := mydb.ExecContext(ctx, ddl); err != nil {
		return this, err
	}

	this.db = mydb
	this.objects_queries = objects.New(this.db)

	return this, nil
}

func sqlxRowsToMapStringAny(query_result *sqlx.Rows) ([]map[string]any, error) {
	var results []map[string]any

	for query_result.Next() {
		res := make(map[string]interface{})
		err := query_result.MapScan(res)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}
