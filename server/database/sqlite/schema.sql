-- user table
CREATE TABLE IF NOT EXISTS users (
    id   INTEGER PRIMARY KEY,
    username text    NOT NULL UNIQUE,
    password_hash BLOB(32) NOT NULL CHECK(length(password_hash) = 32),
    session text NOT NULL UNIQUE
);

-- collections table
CREATE TABLE IF NOT EXISTS collections (
    id   INTEGER PRIMARY KEY,
    table_name text NOT NULL UNIQUE,
    query_rules_directory_path text NOT NULL GENERATED ALWAYS AS ("access_rules/rules/" || id || "/") STORED
);

-- table_fields table 
CREATE TABLE IF NOT EXISTS table_fields (
    id INTEGER PRIMARY KEY,
    collection_id INTEGER NOT NULL,

    field_name text NOT NULL,
    field_type text NOT NULL,
    field_options text,

    is_locked bool NOT NULL DEFAULT FALSE

    is_foreign_key boolean NOT NULL DEFAULT FALSE,
    fk_refers_to_table text,

    FOREIGN KEY(collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

-- querys table
CREATE TABLE IF NOT EXISTS queries (
    id   INTEGER PRIMARY KEY,
    name text UNIQUE NOT NULL,
    query text    NOT NULL DEFAULT 'SELECT * FROM ?',
    query_rules_file_path text NOT NULL GENERATED ALWAYS AS ("access_rules/rules/" || name || ".lua") STORED
);
