package types

const HASH_SIZE = 32
const SESSION_SIZE = 32

type PasswordHash [HASH_SIZE]byte

type Collection struct {
    ID int64
    Name string
    QueryRulesDirectoryPath string
    Fields []TableField
}

type TableField struct {
	ID           int64
	CollectionID int64
	FieldName    string
	FieldType    string
	FieldOptions NullString
}

type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

type Query struct {
	ID                 int64
	Name               string
	Query              string
	QueryRulesFilePath string
}

type User struct {
	ID           int64
	Username     string
	PasswordHash PasswordHash
	Session      string
}

const (
    TOKEN_VARIABLE_TYPE = iota
    TOKEN_NUMBER_TYPE = iota
    TOKEN_OPERATOR = iota
    TOKEN_OPEN_PARENTHESIS = iota
    TOKEN_CLOSE_PARENTHESIS = iota
    TOKEN_EOF = iota
)

var filter_lang_operators = []string{
    "==",
    "!=",
    ">",
    ">=",
    "<",
    "<=",
    "&&",
    "||",
    "~",
    "!~",
    ".",
    "!.",
    ">.",
    ">=.",
    "<.",
    "<=.",
    "~.",
    "!~.",
}

type Token struct {
    Value string
    Type int 
}

