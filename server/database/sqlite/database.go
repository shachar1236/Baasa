package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
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

	logger *slog.Logger
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

    users, err := this.GetCollectionByName(ctx, "users")
    fmt.Printf("%+v\n", users)
    
	if err != nil {
		users, err := this.objects_queries.CreateCollection(ctx, "users")
		if err != nil {
			return this, err
		}
		this.logger.Info("created users collection")

		err = this.objects_queries.CreateField(ctx, objects.CreateFieldParams{
			FieldName:    "id",
			FieldType:    "INTEGER",
			FieldOptions: sql.NullString{String: "PRIMARY KEY", Valid: true},
			CollectionID: users.ID,
			IsLocked:     true,
		})

		if err != nil {
			return this, err
		}

		err = this.objects_queries.CreateField(ctx, objects.CreateFieldParams{
			FieldName:    "username",
			FieldType:    "text",
			FieldOptions: sql.NullString{String: "NOT NULL UNIQUE", Valid: true},
			CollectionID: users.ID,
			IsLocked:     true,
		})

		if err != nil {
			return this, err
		}

		err = this.objects_queries.CreateField(ctx, objects.CreateFieldParams{
			FieldName:    "password_hash",
			FieldType:    "BLOB(32)",
			FieldOptions: sql.NullString{String: "NOT NULL CHECK(length(password_hash) = 32)", Valid: true},
			CollectionID: users.ID,
			IsLocked:     true,
		})

		if err != nil {
			return this, err
		}

		err = this.objects_queries.CreateField(ctx, objects.CreateFieldParams{
			FieldName:    "session",
			FieldType:    "text",
			FieldOptions: sql.NullString{String: "NOT NULL UNIQUE", Valid: true},
			CollectionID: users.ID,
			IsLocked:     true,
		})

		if err != nil {
			return this, err
		}
	}

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
