package database

import (
	"context"

	"github.com/shachar1236/Baasa/database/types"
)

type Database interface {
    GetCollectionByName(ctx context.Context, name string) (types.Collection, error)
    GetCollections(ctx context.Context) ([]types.Collection, error)

    GetQueryById(ctx context.Context, id int64) (string, error)
    GetQuaryById(ctx context.Context, query_id int64) (types.Query, error)
    RunQueryWithFilters(ctx context.Context, query types.Query, args map[string]any, filters string) ([]byte, error)
    RunQuery(query string, args map[string]any) ([]map[string]any, error)
    BasicCount(collection_name string, filter string, args []any) (int64, error)

    DoesUserExists(ctx context.Context, username string, password_hash types.PasswordHash) (bool, error)
    DoesUserExistsById(ctx context.Context, id int64) (bool, error)
    GetUserById(ctx context.Context, id int64) (types.User, error)
    GetUserBySession(ctx context.Context, session string) (types.User, error)
    CreateUser(ctx context.Context, username string, password string) (session string, err error)
    GetUser(ctx context.Context, username string, password string) (types.User, error)
}

