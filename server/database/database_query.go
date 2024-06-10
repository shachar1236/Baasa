package database

import (
	"context"

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

