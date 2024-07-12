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
	ID              int64
	CollectionID    int64
	FieldName       string
	FieldType       string
	FieldOptions    NullString
	IsForeignKey    bool
	FkRefersToTable NullString
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

// ------------- token ----------------


var Filter_lang_operators = []string{
    "=",
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

const (
    TOKEN_VALUE_STRING_TYPE = iota
    TOKEN_VALUE_VARIABLE_TYPE = iota
)

type TokenValue interface {
    GetType() int
}

// type TokenValueVariable struct {
    // Parts []string
    // PartToCollection map[string]Collection
    // PartToCollectionField map[string]TableField
    // UsedCollectionsFilters map[string]string
    // ExtandsToList bool
// }

const (
    TOKEN_VALUE_VARIABLE_PART_COLLECTION_TYPE = iota
    TOKEN_VALUE_VARIABLE_PART_COLLECTION_FIELD_TYPE = iota
)

type TokenValueVariablePartField struct {
    FieldName string
    FieldCollection string
    // if field not fk then fk = ""
    FkRefersToCollection string
}

type TokenValueVariablePartCollection struct {
    CollectionName string
    FkToLastPartName string
}

type TokenValueVariablePart struct {
    PartType int
    Field TokenValueVariablePartField
    Collection TokenValueVariablePartCollection
}

type TokenValueVariable struct {
    Parts []string
    Fields []TokenValueVariablePart
}

func (value TokenValueVariable) GetType() int {
    return TOKEN_VALUE_VARIABLE_TYPE
}

type TokenValueString string

func (value TokenValueString) GetType() int {
    return TOKEN_VALUE_STRING_TYPE
}
const (
    TOKEN_VARIABLE_TYPE = iota
    TOKEN_NUMBER_TYPE = iota
    TOKEN_STRING_TYPE = iota
    TOKEN_OPERATOR = iota
    TOKEN_OPEN_PARENTHESIS = iota
    TOKEN_CLOSE_PARENTHESIS = iota
    TOKEN_EOF = iota
)

type Token struct {
    Value TokenValue
    Type int 
}

