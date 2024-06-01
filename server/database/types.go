package database

const HASH_SIZE = 32
const SESSION_SIZE = 32

type PasswordHash [HASH_SIZE]byte

type TableField struct {
    Name string
    Type string
    Options string
};

type Table struct {
    Name string
    Fields []TableField
};
