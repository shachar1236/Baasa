package database

import (
	"context"

	"github.com/shachar1236/Baasa/database/types"
	querylang_types "github.com/shachar1236/Baasa/query_lang/types"
)

type Database interface {
    GetCollectionByName(ctx context.Context, name string) (types.Collection, error)
    GetCollectionById(ctx context.Context, id int64) (types.Collection, error)
    GetCollections(ctx context.Context) ([]types.Collection, error)
    GetBaseCollections() (collections_names []string , err error)
    AddCollection(ctx context.Context) (types.Collection, error)
    DeleteCollection(ctx context.Context, name string) error
    SaveCollectionChanges(ctx context.Context, new_collection types.Collection) error

    // build and runs user query
    RunUserCustomQuery(collection_name string, fields []string, filter_tokens []querylang_types.Token, sort_by string, expand string, used_collections_filters map[string]string) (resJson string, err error)

    GetQuaries(ctx context.Context) ([]types.Query, error)
    GetQuaryById(ctx context.Context, query_id int64) (types.Query, error)
    CreateQuery(ctx context.Context, name string) (types.Query, error)
    UpdateQuaryById(ctx context.Context, query_id int64, query string) error
    DeleteQuaryById(ctx context.Context, query_id int64) error
    GetQueryRules(ctx context.Context, query_id int64) (string, error)
    SetQueryRules(ctx context.Context, query_id int64, new_rules string) error
    RunQueryWithFilters(ctx context.Context, query types.Query, args map[string]any, filters string) ([]byte, error)
    RunQuery(query string, args map[string]any) ([]map[string]any, error)
    BasicCount(collection_name string, filter string, args []any) (int64, error)
    Get(collection_name string, filters string, args []any) (map[string]any, error)

    ActionSetById(ctx context.Context, collection_name string, column_name string, id int64, to any) error
    ActionGetCollectionDataWithLimit(ctx context.Context, collection_name string, from int64, to int64) ([]map[string]any, error)
    ActionDeleteById(ctx context.Context, collection_name string, id int64) error
    ActionAdd(ctx context.Context, collection_name string, args map[string]any) (int64, error)

    DoesUserExists(ctx context.Context, username string, password_hash types.PasswordHash) (bool, error)
    DoesUserExistsById(ctx context.Context, id int64) (bool, error)
    GetUserById(ctx context.Context, id int64) (types.User, error)
    GetUserBySession(ctx context.Context, session string) (types.User, error)
    CreateUser(ctx context.Context, username string, password string) (session string, err error)
    GetUser(ctx context.Context, username string, password string) (types.User, error)
}

