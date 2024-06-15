package sqlite

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/shachar1236/Baasa/database/types"
)

func (this *SqliteDB) GetQueryById(ctx context.Context, id int64) (string, error) {
    query, err := this.objects_queries.GetQueryById(ctx, id)
    if err != nil {
        this.logger.Error("Cannot get query: ", err)
        return "", err
    }
    return query.Query, nil
}

func (this *SqliteDB) GetQuaryById(ctx context.Context, query_id int64) (types.Query, error) {
    // getting query
    query, err := this.objects_queries.GetQueryById(ctx, query_id)
    return types.Query(query), err
}

// runs query with filters and returns the result as json
// return - json
func (this *SqliteDB) RunQueryWithFilters(ctx context.Context, query types.Query, args map[string]any, filters string) ([]byte, error) {
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
    query_result, err := this.db.NamedQueryContext(ctx, query_string, args)
    if err != nil {
        this.logger.Error("Cannot run query: ", err)
        return nil, err
    }

    var results []map[string]any

    results, err = sqlxRowsToMapStringAny(query_result)
    if err != nil {
        this.logger.Error("Cannot run query: ", err)
        return nil, err
    }

    res, err := json.Marshal(results)
    return res, err
}

func (this *SqliteDB) RunQuery(query string, args map[string]any) ([]map[string]any, error) {
    // running query
    query_result, err := this.db.NamedQuery(query, args)
    if err != nil {
        this.logger.Error("Cannot run query: ", err)
        return nil, err
    }

    var results []map[string]any

    for query_result.Next() {
        res := make(map[string]interface{})
        err := query_result.MapScan(res)
        if err != nil {
            this.logger.Error("Cannot run query: ", err)
            return nil, err
        }
        results = append(results, res)
    }

    return results, nil
}

func (this *SqliteDB) BasicCount(collection_name string, filters string, args []any) (int64, error) {
    query_select := "SELECT COUNT(*) FROM "
    query_where := " WHERE "

    var sb strings.Builder
    sb.Grow(len(query_select) + len(collection_name) + len(query_where) + len(filters) + 1)
    sb.WriteString(query_select)
    sb.WriteString(collection_name)
    sb.WriteString(query_where)
    sb.WriteString(filters)
    sb.WriteString(";")

    var count int64
    // running query
    err := this.db.Get(&count, sb.String(), args...)
    return count, err
}
