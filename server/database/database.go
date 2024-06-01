package database

import (
	"context"
	"crypto/sha256"
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
var my_tables []Table

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
    my_tables = append(my_tables, Table{
        Name: "users",
        Fields: []TableField{
            TableField{
                Name: "id",
                Type: "INTEGER",
                Options: "PRIMARY KEY",
            },
            TableField{
                Name: "username",
                Type: "TEXT",
                Options: "NOT NULL UNIQUE",
            },
            TableField{
                Name: "password_hash",
                Type: "BLOB(32)",
                Options: "NOT NULL CHECK(length(password_hash) = 32)",
            },
            TableField{
                Name: "session",
                Type: "TEXT",
                Options: "NOT NULL UNIQUE",
            },
        },
    })

    // create user tables
    tables, err := GetTables(ctx)
    if err != nil {
        panic(err)
    }

    if len(tables) > 0 {
        var curr_table Table
        for _, table := range tables {
            if table.TableName != curr_table.Name {
                // add last table
                my_tables = append(my_tables, curr_table)

                // new table
                curr_table = Table{}
                curr_table.Name = table.TableName
            } 
            curr_table.Fields = append(curr_table.Fields, TableField{
                Name: table.FieldName,
                Type: table.FieldName,
                Options: table.FieldOptions.String,
            })
        }
    }
}

func DoesUserExists(ctx context.Context, username string, password_hash PasswordHash) (bool, error) {
    count, err := objects_queries.CountUsersWithNameAndPassword(ctx, 
        objects.CountUsersWithNameAndPasswordParams{ Username: username, PasswordHash: password_hash[:]})
    return count > 0, err
}

func DoesUserExistsById(ctx context.Context, id int64) (bool, error) {
    count, err := objects_queries.CountUsersWithId(ctx, id)
    return count > 0, err
}

func GetUserById(ctx context.Context, id int64) (objects.User, error) {
    user, err := objects_queries.GetUserById(ctx, id)
    return user, err
}

func GetUserBySession(ctx context.Context, session string) (objects.User, error) {
    user, err := objects_queries.GetUserBySession(ctx, session)
    return user, err
}

func CreateUser(ctx context.Context, username string, password string) (session string, err error){
    // generating password hash
    password_hash := sha256.Sum256([]byte(password))
    // generate random session
    rand_session, err := GenerateSecureToken()
    if err != nil {
        return "", err
    }
    // creating user
    err = objects_queries.CreateUser(ctx, 
        objects.CreateUserParams{Username: username, PasswordHash: password_hash[:], Session: rand_session})
    return rand_session, err
}

func GetUser(ctx context.Context, username string, password string) (objects.User, error) {
    password_hash := sha256.Sum256([]byte(password))
    user, err := objects_queries.GetUserByNameAndPassword(ctx, 
        objects.GetUserByNameAndPasswordParams{Username: username, PasswordHash: password_hash[:]}) 
    return user, err
}

func GetQueryById(ctx context.Context, id int64) (string, error) {
    query, err := objects_queries.GetQueryById(ctx, id)
    if err != nil {
        return "", err
    }
    return query.Query, nil
}

func GetTableByName(ctx context.Context, name string) (objects.GetTableAndFielsByTableNameRow, error) {
    res, err := objects_queries.GetTableAndFielsByTableName(ctx, name)
    return res, err
}

func GetTables(ctx context.Context) ([]objects.GetAllTablesAndFieldsRow, error) {
    res, err := objects_queries.GetAllTablesAndFields(ctx)
    return res, err
}

func GetQuaryById(ctx context.Context, query_id int64) (objects.Query, error) {
    // getting query
    query, err := objects_queries.GetQueryById(ctx, query_id)
    return query, err
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
