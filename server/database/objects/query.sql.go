// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package objects

import (
	"context"
	"database/sql"
)

const countUsersWithId = `-- name: CountUsersWithId :one
SELECT count(*) 
FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) CountUsersWithId(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsersWithId, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countUsersWithNameAndPassword = `-- name: CountUsersWithNameAndPassword :one
SELECT count(*) 
FROM users
WHERE username = ? and password_hash = ? LIMIT 1
`

type CountUsersWithNameAndPasswordParams struct {
	Username     string
	PasswordHash interface{}
}

func (q *Queries) CountUsersWithNameAndPassword(ctx context.Context, arg CountUsersWithNameAndPasswordParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsersWithNameAndPassword, arg.Username, arg.PasswordHash)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createField = `-- name: CreateField :exec
INSERT INTO table_fields (field_name, field_type, field_options, collection_id)
VALUES (?, ?, ?, ?)
`

type CreateFieldParams struct {
	FieldName    string
	FieldType    string
	FieldOptions sql.NullString
	CollectionID int64
}

func (q *Queries) CreateField(ctx context.Context, arg CreateFieldParams) error {
	_, err := q.db.ExecContext(ctx, createField,
		arg.FieldName,
		arg.FieldType,
		arg.FieldOptions,
		arg.CollectionID,
	)
	return err
}

const createQuery = `-- name: CreateQuery :exec
INSERT INTO queries (query)
VALUES (?)
`

func (q *Queries) CreateQuery(ctx context.Context, query string) error {
	_, err := q.db.ExecContext(ctx, createQuery, query)
	return err
}

const createTable = `-- name: CreateTable :exec
INSERT INTO collections (table_name)
VALUES (?) RETURNING id, table_name, query_rules_directory_path
`

func (q *Queries) CreateTable(ctx context.Context, tableName string) error {
	_, err := q.db.ExecContext(ctx, createTable, tableName)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (username, password_hash, session)
VALUES (?, ?, ?)
`

type CreateUserParams struct {
	Username     string
	PasswordHash interface{}
	Session      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.Username, arg.PasswordHash, arg.Session)
	return err
}

const deleteQueryById = `-- name: DeleteQueryById :exec
DELETE FROM queries
WHERE id = ?
`

func (q *Queries) DeleteQueryById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteQueryById, id)
	return err
}

const getAllTablesAndFields = `-- name: GetAllTablesAndFields :many
SELECT collections.id AS collection_id, 
       collections.table_name, 
       collections.query_rules_directory_path as QueryRulesDirectoryPath,
       table_fields.id AS field_id, 
       table_fields.field_name, 
       table_fields.field_type, 
       table_fields.field_options
FROM collections 
INNER JOIN table_fields ON table_fields.collection_id = collections.id
`

type GetAllTablesAndFieldsRow struct {
	CollectionID            int64
	TableName               string
	Queryrulesdirectorypath string
	FieldID                 int64
	FieldName               string
	FieldType               string
	FieldOptions            sql.NullString
}

func (q *Queries) GetAllTablesAndFields(ctx context.Context) ([]GetAllTablesAndFieldsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllTablesAndFields)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllTablesAndFieldsRow
	for rows.Next() {
		var i GetAllTablesAndFieldsRow
		if err := rows.Scan(
			&i.CollectionID,
			&i.TableName,
			&i.Queryrulesdirectorypath,
			&i.FieldID,
			&i.FieldName,
			&i.FieldType,
			&i.FieldOptions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueryById = `-- name: GetQueryById :one
SELECT id, "query", query_rules_file_path FROM queries
WHERE id = ? LIMIT 1
`

func (q *Queries) GetQueryById(ctx context.Context, id int64) (Query, error) {
	row := q.db.QueryRowContext(ctx, getQueryById, id)
	var i Query
	err := row.Scan(&i.ID, &i.Query, &i.QueryRulesFilePath)
	return i, err
}

const getQueryByQuery = `-- name: GetQueryByQuery :one
SELECT id, "query", query_rules_file_path FROM queries
WHERE query = ? LIMIT 1
`

func (q *Queries) GetQueryByQuery(ctx context.Context, query string) (Query, error) {
	row := q.db.QueryRowContext(ctx, getQueryByQuery, query)
	var i Query
	err := row.Scan(&i.ID, &i.Query, &i.QueryRulesFilePath)
	return i, err
}

const getTableAndFielsByTableName = `-- name: GetTableAndFielsByTableName :one
SELECT collections.id, table_name, query_rules_directory_path, table_fields.id, collection_id, field_name, field_type, field_options FROM collections
INNER JOIN table_fields
ON table_fields.collection_id = collections.id
WHERE collections.table_name = ?
`

type GetTableAndFielsByTableNameRow struct {
	ID                      int64
	TableName               string
	QueryRulesDirectoryPath string
	ID_2                    int64
	CollectionID            int64
	FieldName               string
	FieldType               string
	FieldOptions            sql.NullString
}

func (q *Queries) GetTableAndFielsByTableName(ctx context.Context, tableName string) (GetTableAndFielsByTableNameRow, error) {
	row := q.db.QueryRowContext(ctx, getTableAndFielsByTableName, tableName)
	var i GetTableAndFielsByTableNameRow
	err := row.Scan(
		&i.ID,
		&i.TableName,
		&i.QueryRulesDirectoryPath,
		&i.ID_2,
		&i.CollectionID,
		&i.FieldName,
		&i.FieldType,
		&i.FieldOptions,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, password_hash, session FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.Session,
	)
	return i, err
}

const getUserByNameAndPassword = `-- name: GetUserByNameAndPassword :one
SELECT id, username, password_hash, session FROM users
WHERE username = ? and password_hash = ? LIMIT 1
`

type GetUserByNameAndPasswordParams struct {
	Username     string
	PasswordHash interface{}
}

func (q *Queries) GetUserByNameAndPassword(ctx context.Context, arg GetUserByNameAndPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByNameAndPassword, arg.Username, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.Session,
	)
	return i, err
}

const getUserBySession = `-- name: GetUserBySession :one
SELECT id, username, password_hash, session FROM users
WHERE session = ? LIMIT 1
`

func (q *Queries) GetUserBySession(ctx context.Context, session string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserBySession, session)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.Session,
	)
	return i, err
}

const updateQueryById = `-- name: UpdateQueryById :exec
UPDATE queries
SET query = ?
WHERE id = ?
`

type UpdateQueryByIdParams struct {
	Query string
	ID    int64
}

func (q *Queries) UpdateQueryById(ctx context.Context, arg UpdateQueryByIdParams) error {
	_, err := q.db.ExecContext(ctx, updateQueryById, arg.Query, arg.ID)
	return err
}
