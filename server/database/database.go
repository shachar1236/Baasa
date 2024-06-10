package database

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"
	"strings"

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
                FieldOptions: sql.NullString{String:"PRIMARY KEY",, Valid: true},
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
    tables, err := GetCollections(ctx)
    if err != nil {
        panic(err)
    }
}

// runs query with filters and returns the result as json
// return - json
func RunQueryWithFilters(ctx context.Context, query objects.Query, args map[string]any, filters string) (string, error) {


    var query_string string
    
    if filters != "" {
        select_start := "SELECT * FROM ("
        select_end := ") WHERE "
        var sb strings.Builder
        sb.Grow(len(select_start) + len(query.Query) + len(select_end) + len(filters) + 1)

        sb.WriteString(select_start)
        if query.Query[len(query.Query)-1] == ';' {
            sb.WriteString(query.Query[:len(query.Query)-1])
        } else {
            sb.WriteString(query.Query)
        }
        sb.WriteString(select_end)
        sb.WriteString(filters)
        sb.WriteString(";")

        query_string = sb.String()
        log.Println(query_string)
    } else {
        query_string = query.Query
    }

    // running query
    query_result, err := db.NamedQueryContext(ctx, query_string, args)
    if err != nil {
        return "", err
    }

    var results []map[string]any

    for query_result.Next() {
        res := make(map[string]interface{})
        err := query_result.MapScan(res)
        if err != nil {
            return "", err
        }
        results = append(results, res)
    }

    res, err := json.Marshal(results)
    return string(res), err
}

func RunQuery(query string, args map[string]any) ([]map[string]any, error) {
    // running query
    query_result, err := db.NamedQuery(query, args)
    if err != nil {
        return nil, err
    }

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

func RunCountQuery(query string, args []any) (int64, error) {
    var count int64
    // running query
    err := db.Get(&count, query, args...)
    return count, err
}
