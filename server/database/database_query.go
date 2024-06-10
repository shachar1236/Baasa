package database

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/shachar1236/Baasa/database/objects"
)

func GetQueryById(ctx context.Context, id int64) (string, error) {
    query, err := objects_queries.GetQueryById(ctx, id)
    if err != nil {
        return "", err
    }
    return query.Query, nil
}

func GetQuaryById(ctx context.Context, query_id int64) (objects.Query, error) {
    // getting query
    query, err := objects_queries.GetQueryById(ctx, query_id)
    return query, err
}

// runs query with filters and returns the result as json
// return - json
func RunQueryWithFilters(ctx context.Context, query objects.Query, args map[string]any, filters string) ([]byte, error) {
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

    res, err := json.Marshal(results)
    return res, err
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
