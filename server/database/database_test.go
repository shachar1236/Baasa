package database

import (
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"testing"
)

func TestQuary(t *testing.T) {
    var err error
	db, err = sqlx.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

    quary_result, err  := db.Queryx("SELECT * FROM my_tables")
    if err != nil {
        t.Error(err)
    }

    var results []map[string]any

    for quary_result.Next() {
        res := make(map[string]interface{})
        quary_result.MapScan(res)
        results = append(results, res)
    }

    res, err := json.Marshal(results)
    if err != nil {
        t.Error(err)
    }
    t.Log(string(res))
    t.Error(string("hi"))
}
