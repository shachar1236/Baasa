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
	"github.com/shachar1236/Baasa/database/types"
)

//go:embed schema.sql
var ddl string

type SqliteDB struct {
    db *sqlx.DB
    objects_queries *objects.Queries

    // tables
    my_tables []types.Collection
    logger *slog.Logger
}

func New(ctx context.Context) (SqliteDB, error) {
    logFile, err := os.OpenFile("logs/sqlite_db.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
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

    // add users table to my_tables
    this.my_tables = append(this.my_tables, types.Collection{
        Name: "users",
        Fields: []types.TableField{
            types.TableField{
                FieldName: "id",
                FieldType: "INTEGER",
                FieldOptions: types.NullString{String:"PRIMARY KEY", Valid: true},
            },
            types.TableField{
                FieldName: "username",
                FieldType: "TEXT",
                FieldOptions: types.NullString{String: "NOT NULL UNIQUE", Valid: true},
            },
            types.TableField{
                FieldName: "password_hash",
                FieldType: "BLOB(32)",
                FieldOptions: types.NullString{String:"NOT NULL CHECK(length(password_hash) = 32)", Valid: true},
            },
            types.TableField{
                FieldName: "session",
                FieldType: "TEXT",
                FieldOptions: types.NullString{String:"NOT NULL UNIQUE", Valid: true},
            },
        },
    })

    // create user tables
    // tables, err := GetCollections(ctx)
    if err != nil {
        return this, err
    }

    return this, nil
}

