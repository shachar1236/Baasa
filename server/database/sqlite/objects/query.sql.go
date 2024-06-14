// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package objects

import (
	"context"
	"database/sql"
)

const changeFieldName = `-- name: ChangeFieldName :exec
UPDATE table_fields SET field_name = ? WHERE id = ?
`

type ChangeFieldNameParams struct {
	FieldName string
	ID        int64
}

func (q *Queries) ChangeFieldName(ctx context.Context, arg ChangeFieldNameParams) error {
	_, err := q.db.ExecContext(ctx, changeFieldName, arg.FieldName, arg.ID)
	return err
}

const changeFieldOptions = `-- name: ChangeFieldOptions :exec
UPDATE table_fields SET field_options = ? WHERE id = ?
`

type ChangeFieldOptionsParams struct {
	FieldOptions sql.NullString
	ID           int64
}

func (q *Queries) ChangeFieldOptions(ctx context.Context, arg ChangeFieldOptionsParams) error {
	_, err := q.db.ExecContext(ctx, changeFieldOptions, arg.FieldOptions, arg.ID)
	return err
}

const changeFieldType = `-- name: ChangeFieldType :exec
UPDATE table_fields SET field_type = ? WHERE id = ?
`

type ChangeFieldTypeParams struct {
	FieldType string
	ID        int64
}

func (q *Queries) ChangeFieldType(ctx context.Context, arg ChangeFieldTypeParams) error {
	_, err := q.db.ExecContext(ctx, changeFieldType, arg.FieldType, arg.ID)
	return err
}

const changeTableName = `-- name: ChangeTableName :exec
UPDATE collections SET table_name = ?1 WHERE table_name = ?2
`

type ChangeTableNameParams struct {
	NewName string
	OldName string
}

func (q *Queries) ChangeTableName(ctx context.Context, arg ChangeTableNameParams) error {
	_, err := q.db.ExecContext(ctx, changeTableName, arg.NewName, arg.OldName)
	return err
}

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

const createCollection = `-- name: CreateCollection :one
INSERT INTO collections (table_name)
VALUES (?) RETURNING id, table_name, query_rules_directory_path
`

func (q *Queries) CreateCollection(ctx context.Context, tableName string) (Collection, error) {
	row := q.db.QueryRowContext(ctx, createCollection, tableName)
	var i Collection
	err := row.Scan(&i.ID, &i.TableName, &i.QueryRulesDirectoryPath)
	return i, err
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

const deleteCollection = `-- name: DeleteCollection :exec
DELETE FROM collections WHERE table_name = ?
`

func (q *Queries) DeleteCollection(ctx context.Context, tableName string) error {
	_, err := q.db.ExecContext(ctx, deleteCollection, tableName)
	return err
}

const deleteField = `-- name: DeleteField :exec
DELETE FROM table_fields
WHERE id = ?
`

func (q *Queries) DeleteField(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteField, id)
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
LEFT JOIN table_fields ON collections.id = table_fields.collection_id
`

type GetAllTablesAndFieldsRow struct {
	CollectionID            int64
	TableName               string
	Queryrulesdirectorypath string
	FieldID                 sql.NullInt64
	FieldName               sql.NullString
	FieldType               sql.NullString
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

const getTableAndFieldsByTableId = `-- name: GetTableAndFieldsByTableId :many
SELECT collections.id AS collection_id,
    collections.table_name,
    collections.query_rules_directory_path as QueryRulesDirectoryPath,
    table_fields.id AS field_id,
    table_fields.field_name,
    table_fields.field_type,
    table_fields.field_options
FROM collections
    LEFT JOIN table_fields ON collections.id = table_fields.collection_id
WHERE collections.id = ?
`

type GetTableAndFieldsByTableIdRow struct {
	CollectionID            int64
	TableName               string
	Queryrulesdirectorypath string
	FieldID                 sql.NullInt64
	FieldName               sql.NullString
	FieldType               sql.NullString
	FieldOptions            sql.NullString
}

func (q *Queries) GetTableAndFieldsByTableId(ctx context.Context, id int64) ([]GetTableAndFieldsByTableIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getTableAndFieldsByTableId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTableAndFieldsByTableIdRow
	for rows.Next() {
		var i GetTableAndFieldsByTableIdRow
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

const getTableAndFieldsByTableName = `-- name: GetTableAndFieldsByTableName :many
SELECT collections.id AS collection_id,
    collections.table_name,
    collections.query_rules_directory_path as QueryRulesDirectoryPath,
    table_fields.id AS field_id,
    table_fields.field_name,
    table_fields.field_type,
    table_fields.field_options
FROM collections
    LEFT JOIN table_fields ON collections.id = table_fields.collection_id WHERE collections.table_name = ?
`

type GetTableAndFieldsByTableNameRow struct {
	CollectionID            int64
	TableName               string
	Queryrulesdirectorypath string
	FieldID                 sql.NullInt64
	FieldName               sql.NullString
	FieldType               sql.NullString
	FieldOptions            sql.NullString
}

func (q *Queries) GetTableAndFieldsByTableName(ctx context.Context, tableName string) ([]GetTableAndFieldsByTableNameRow, error) {
	rows, err := q.db.QueryContext(ctx, getTableAndFieldsByTableName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTableAndFieldsByTableNameRow
	for rows.Next() {
		var i GetTableAndFieldsByTableNameRow
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
