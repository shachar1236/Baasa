package sqlite

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/shachar1236/Baasa/database/sqlite/objects"
	"github.com/shachar1236/Baasa/database/types"
)

func (this *SqliteDB) CreateQuery(ctx context.Context, name string) (types.Query, error) {
    res, err := this.objects_queries.CreateQuery(ctx, name)
    return types.Query(res), err
}

func (this *SqliteDB) GetQuaries(ctx context.Context) ([]types.Query, error) {
    res, err := this.objects_queries.GetQueries(ctx)
    if err != nil {
        this.logger.Error("Cannot get queries: ", err)
        return nil, err
    }
    ret := make([]types.Query, len(res))
    for i, q := range res {
        ret[i] = types.Query(q)
    }

    return ret, err
}

func (this *SqliteDB) UpdateQuaryById(ctx context.Context, query_id int64, query string) error {
    err := this.objects_queries.UpdateQueryById(ctx, objects.UpdateQueryByIdParams{
        Query: query,
        ID: query_id,
    })
    return err
}

func (this *SqliteDB) DeleteQuaryById(ctx context.Context, query_id int64) error {
    // getting query
    query, err := this.objects_queries.GetQueryById(ctx, query_id)
    if err != nil {
        this.logger.Error("Cannot get query: ", err)
        return err
    }
    // removing query file
    err = os.Remove(query.QueryRulesFilePath)
    if err != nil {
        this.logger.Error("Cannot remove query file: ", err)
        return err
    }
    err = this.objects_queries.DeleteQueryById(ctx, query_id)
    return err
}

func (this *SqliteDB) GetQuaryById(ctx context.Context, query_id int64) (types.Query, error) {
    // getting query
    query, err := this.objects_queries.GetQueryById(ctx, query_id)
    return types.Query(query), err
}

func (this *SqliteDB) GetQueryRules(ctx context.Context, query_id int64) (string, error) {
    query, err := this.GetQuaryById(ctx, query_id)
    if err != nil {
        this.logger.Error("Cannot get query: ", err)
        return "", err
    }

    file, err := os.Open(query.QueryRulesFilePath)
    
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            file, err = os.Create(query.QueryRulesFilePath)
            if err != nil {
                this.logger.Error("Cannot create file: ", err)
                return "", err
            }
            return "", nil
        } else {
            this.logger.Error("Cannot open file: ", err)
            return "", err
        }
    }
    defer file.Close()

    content, err := io.ReadAll(file)
    
    return string(content), err
}

func (this *SqliteDB) SetQueryRules(ctx context.Context, query_id int64, new_rules string) error {
    query, err := this.GetQuaryById(ctx, query_id)
    if err != nil {
        this.logger.Error("Cannot get query: ", err)
        return err
    }

    file, err := os.OpenFile(query.QueryRulesFilePath, os.O_WRONLY, 0644)
    
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            file, err = os.Create(query.QueryRulesFilePath)
            if err != nil {
                this.logger.Error("Cannot create file: ", err)
                return err
            }
        } else {
            this.logger.Error("Cannot open file: ", err)
            return err
        }
    }
    defer file.Close()

    // remove all contents of file and write new rules
    err = file.Truncate(0)
    if err != nil {
        this.logger.Error("Cannot truncate file: ", err)
        return err
    }
    _, err = file.WriteString(new_rules)
    
    return err
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
