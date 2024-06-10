-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserBySession :one
SELECT * FROM users
WHERE session = ? LIMIT 1;

-- name: GetUserByNameAndPassword :one
SELECT * FROM users
WHERE username = ? and password_hash = ? LIMIT 1;

-- name: CountUsersWithNameAndPassword :one
SELECT count(*) 
FROM users
WHERE username = ? and password_hash = ? LIMIT 1;

-- name: CountUsersWithId :one
SELECT count(*) 
FROM users
WHERE id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (username, password_hash, session)
VALUES (?, ?, ?);

-- name: CreateQuery :exec
INSERT INTO queries (query)
VALUES (?);

-- name: GetQueryById :one
SELECT * FROM queries
WHERE id = ? LIMIT 1;

-- name: GetQueryByQuery :one
SELECT * FROM queries
WHERE query = ? LIMIT 1;

-- name: DeleteQueryById :exec
DELETE FROM queries
WHERE id = ?;

-- name: UpdateQueryById :exec
UPDATE queries
SET query = ?
WHERE id = ?;

-- name: GetTableAndFielsByTableName :one 
SELECT * FROM collections
INNER JOIN table_fields
ON table_fields.collection_id = collections.id
WHERE collections.table_name = ?;

-- name: GetAllTablesAndFields :many
SELECT collections.id AS collection_id, 
       collections.table_name, 
       collections.query_rules_directory_path as QueryRulesDirectoryPath,
       table_fields.id AS field_id, 
       table_fields.field_name, 
       table_fields.field_type, 
       table_fields.field_options
FROM collections 
INNER JOIN table_fields ON table_fields.collection_id = collections.id;

-- name: CreateTable :exec
INSERT INTO collections (table_name)
VALUES (?) RETURNING *;

-- name: CreateField :exec
INSERT INTO table_fields (field_name, field_type, field_options, collection_id)
VALUES (?, ?, ?, ?);
