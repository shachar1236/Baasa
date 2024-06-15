package sqlite

import (
	"context"
    "golang.org/x/exp/maps"
	"strings"
)

func (this *SqliteDB) ActionSetById(ctx context.Context, collection_name string, column_name string, id int64, to any) error {
    var sb strings.Builder
    base_1 := "UPDATE "
    base_2 := " SET "
    base_3 := " = ? WHERE id = ?;"
    sb.Grow(len(base_1) + len(collection_name) + len(base_2) + len(column_name) + len(base_3))
    sb.WriteString(base_1) // UPDATE
    sb.WriteString(collection_name) // TableName
    sb.WriteString(base_2) // SET
    sb.WriteString(column_name) // ColumnName
    sb.WriteString(base_3) // = ? WHERE id = ?;

    _, err := this.db.Exec(sb.String(), to, id)
    return err
}

func (this *SqliteDB) ActionGetCollectionDataWithLimit(ctx context.Context, collection_name string, from int64, to int64) ([]map[string]any, error) {
    var sb strings.Builder
    base_1 := "SELECT * FROM "
    base_2 := " LIMIT ? OFFSET ?;"
    sb.Grow(len(base_1) + len(collection_name) + len(base_2))
    sb.WriteString(base_1)
    sb.WriteString(collection_name)
    sb.WriteString(base_2)

    logger := this.logger.With("query", sb.String())

    rows, err := this.db.QueryxContext(ctx, sb.String(), to - from, from)
    if err != nil {
        logger.Error("Cannot run query: " + err.Error())
        return nil, err
    }

    var results []map[string]any

    results, err = sqlxRowsToMapStringAny(rows)
    if err != nil {
        logger.Error("Cannot run query: " + err.Error())
        return nil, err
    }

    return results, nil
}

func (this *SqliteDB) ActionDeleteById(ctx context.Context, collection_name string, id int64) error {
    var sb strings.Builder
    base_1 := "DELETE FROM "
    base_2 := " WHERE id = ?;"
    sb.Grow(len(base_1) + len(collection_name) + len(base_2))
    sb.WriteString(base_1) // DELETE FROM
    sb.WriteString(collection_name) // TableName
    sb.WriteString(base_2) // WHERE id = ?;

    _, err := this.db.Exec(sb.String(), id)
    return err
}

func (this *SqliteDB) ActionAdd(ctx context.Context, collection_name string, args map[string]any) (int64, error) {
    var sb strings.Builder
    base_1 := "INSERT INTO "
    base_2 := " ("
    base_3 := ") VALUES ("
    base_4 := ") RETURNING id;"

    args_keys := maps.Keys(args)
    keys := strings.Join(args_keys, ", ")
    sb.Grow(len(base_1) + len(collection_name) + len(base_2) + len(keys) + len(base_3) + len(args_keys) * 2 + len(base_4))
    sb.WriteString(base_1)
    sb.WriteString(collection_name)
    sb.WriteString(base_2)
    sb.WriteString(keys)
    sb.WriteString(base_3)
    sb.WriteString("?")
    for i := 0; i < len(args_keys) - 1; i++ {
        sb.WriteString(",?")
    }
    sb.WriteString(base_4)

    var id int64
    err := this.db.Get(&id, sb.String(), maps.Values(args)...)

    if err != nil {
        this.logger.Error("Cannot run query: " + err.Error() + "; query: " + sb.String())
    }

    return id, err
}
