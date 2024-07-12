package querylang_types

import "github.com/shachar1236/Baasa/database/types"

type CollectionsSet struct {
	collections map[string]types.Collection
}

func (set *CollectionsSet) Init() {
	set.collections = make(map[string]types.Collection)
}

func (set *CollectionsSet) Add(collection types.Collection) {
	set.collections[collection.Name] = collection
}

func (set *CollectionsSet) GetMap() map[string]types.Collection {
	return set.collections
}

func (set *CollectionsSet) Union(added_collection_set CollectionsSet) {
	for k, v := range added_collection_set.collections {
		set.collections[k] = v
	}
}

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

const (
    TOKEN_VALUE_STRING_TYPE = iota
    TOKEN_VALUE_VARIABLE_TYPE = iota
)

type TokenValue interface {
    GetType() int
}

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
