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
    query_rules_directory_path text NOT NULL
);

-- table_fields table 
CREATE TABLE IF NOT EXISTS table_fields (
    id INTEGER PRIMARY KEY,
    collection_id INTEGER NOT NULL,

    field_name text NOT NULL,
    field_type text NOT NULL,
    field_options text,

    FOREIGN KEY(collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

-- querys table
CREATE TABLE IF NOT EXISTS queries (
    id   INTEGER PRIMARY KEY,
    query text    NOT NULL,
    query_rules_file_path text NOT NULL
);
