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

-- name: CreateQuery :one
INSERT INTO queries (name)
VALUES (?) RETURNING *;

-- name: GetQueryById :one
SELECT * FROM queries
WHERE id = ? LIMIT 1;

-- name: GetQueries :many
SELECT * FROM queries;

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

-- name: GetTableAndFieldsByTableName :many 
SELECT collections.id AS collection_id,
    collections.table_name,
    collections.query_rules_directory_path as QueryRulesDirectoryPath,
    table_fields.id AS field_id,
    table_fields.field_name,
    table_fields.field_type,
    table_fields.field_options,
    table_fields.is_foreign_key,
    table_fields.fk_table_name,
    table_fields.fk_field_name
FROM collections
    LEFT JOIN table_fields ON collections.id = table_fields.collection_id WHERE collections.table_name = ?;

-- name: GetTableAndFieldsByTableId :many 
SELECT collections.id AS collection_id,
    collections.table_name,
    collections.query_rules_directory_path as QueryRulesDirectoryPath,
    table_fields.id AS field_id,
    table_fields.field_name,
    table_fields.field_type,
    table_fields.field_options,
    table_fields.is_foreign_key,
    table_fields.fk_table_name,
    table_fields.fk_field_name
FROM collections
    LEFT JOIN table_fields ON collections.id = table_fields.collection_id
WHERE collections.id = ?;

-- name: GetAllTablesAndFields :many
SELECT collections.id AS collection_id, 
       collections.table_name, 
       collections.query_rules_directory_path as QueryRulesDirectoryPath,
       table_fields.id AS field_id, 
       table_fields.field_name, 
       table_fields.field_type, 
       table_fields.field_options,
       table_fields.is_foreign_key,
       table_fields.fk_table_name,
       table_fields.fk_field_name
FROM collections 
LEFT JOIN table_fields ON collections.id = table_fields.collection_id;

-- name: CreateCollection :one
INSERT INTO collections (table_name)
VALUES (?) RETURNING *;

-- name: CreateField :exec
INSERT INTO table_fields (field_name, field_type, field_options, collection_id, is_foreign_key, fk_table_name, fk_field_name)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: DeleteCollection :exec
DELETE FROM collections WHERE table_name = ?;

-- name: ChangeTableName :exec
UPDATE collections SET table_name = sqlc.arg(new_name) WHERE table_name = sqlc.arg(old_name);

-- name: ChangeFieldName :exec
UPDATE table_fields SET field_name = ? WHERE id = ?;

-- name: ChangeFieldType :exec
UPDATE table_fields SET field_type = ? WHERE id = ?;

-- name: ChangeFieldOptions :exec
UPDATE table_fields SET field_options = ? WHERE id = ?;

-- name: ChangeFieldToForeignKey :exec
UPDATE table_fields SET field_type = ?, is_foreign_key = true, fk_table_name = ?, fk_field_name = ? WHERE id = ?;

-- name: ChangeFieldToNotBeForeignKey :exec
UPDATE table_fields SET is_foreign_key = false, fk_table_name = null, fk_field_name = null WHERE id = ?;


-- name: DeleteField :exec
DELETE FROM table_fields
WHERE id = ?;
