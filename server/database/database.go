package database

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shachar1236/Baasa/database/objects"
)

//go:embed schema.sql
var ddl string

var db *sqlx.DB
var objects_queries *objects.Queries

// tables
var my_tables []Collection

func Init(ctx context.Context) {
	mydb, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := mydb.ExecContext(ctx, ddl); err != nil {
        panic(err)
	}

    db = mydb
    objects_queries = objects.New(db)

    // add users table to my_tables
    my_tables = append(my_tables, Collection{
        Name: "users",
        Fields: []objects.TableField{
            objects.TableField{
                FieldName: "id",
                FieldType: "INTEGER",
                FieldOptions: sql.NullString{String:"PRIMARY KEY", Valid: true},
            },
            objects.TableField{
                FieldName: "username",
                FieldType: "TEXT",
                FieldOptions: sql.NullString{String: "NOT NULL UNIQUE", Valid: true},
            },
            objects.TableField{
                FieldName: "password_hash",
                FieldType: "BLOB(32)",
                FieldOptions: sql.NullString{String:"NOT NULL CHECK(length(password_hash) = 32)", Valid: true},
            },
            objects.TableField{
                FieldName: "session",
                FieldType: "TEXT",
                FieldOptions: sql.NullString{String:"NOT NULL UNIQUE", Valid: true},
            },
        },
    })

    // create user tables
    // tables, err := GetCollections(ctx)
    if err != nil {
        panic(err)
    }
}

