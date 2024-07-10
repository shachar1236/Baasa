package sqlite

import (
	"fmt"
	"strings"
	"sync"

	"github.com/leporo/sqlf"
	"github.com/shachar1236/Baasa/database/types"
)

var sb = strings.Builder{}
var sb_mutex = sync.Mutex{}

// build and runs user query
func (db *SqliteDB) RunUserCustomQuery(
    collection_name string,
    fields []string,
    filter_tokens []types.Token,
    sort_by string,
    expand string,
    used_collections_filters []string,
) (resJson string, err error) {
    _ = sqlf.Select(strings.Join(fields, ",")).From(collection_name)

    sb_mutex.Lock()
    defer sb_mutex.Unlock()
    defer sb.Reset()

    for _, token := range filter_tokens {
        if token.Type == types.TOKEN_OPERATOR {
            // TODO: change for supporting special operators
            token_as_string := token.Value.(types.TokenValueString)
            sb.WriteString(string(token_as_string))
        } else if token.Type == types.TOKEN_VARIABLE_TYPE {
            token_as_variable := token.Value.(types.TokenValueVariable)
            variable_parts := token_as_variable.Parts
            if len(variable_parts) > 2 {
                // its nested collections
                createExpandedSelect(&token_as_variable, &sb)
            } else if len(variable_parts) == 2 {
                // its a field from the collection
                sb.WriteString(token_as_variable.Parts[0] + "." + token_as_variable.Parts[1])
            } else {
                // not valid
            }
        } else {
            token_as_string := token.Value.(types.TokenValueString)
            sb.WriteString(string(token_as_string))
        }
        sb.WriteString(" ")
    }

    fmt.Println("Printing user custom query")
    fmt.Println()
    fmt.Println()
    fmt.Println(sb.String())
    fmt.Println()
    fmt.Println()

    return
}

func createExpandedSelect(token_as_variable *types.TokenValueVariable, sb *strings.Builder) {
    // TODO: add collection filters
    variable_parts := token_as_variable.Parts
    last_index := len(variable_parts) - 1
    sb.WriteString("(SELECT ")
    sb.WriteString(variable_parts[last_index])
    sb.WriteString(" FROM ")
    for i := last_index; i > 1; i-- {
        sb.WriteString(" (SELECT ")
        sb.WriteString(variable_parts[i])
        sb.WriteString(" FROM ")
        sb.WriteString(token_as_variable.PartToCollection[variable_parts[i - 1]].Name)
        sb.WriteString(" WHERE ")
        sb.WriteString(token_as_variable.PartToCollectionField[variable_parts[i - 1]].FkFieldName.String)
        sb.WriteString(" = ")
    }
    sb.WriteString(variable_parts[0])
    sb.WriteString(".")
    sb.WriteString(variable_parts[1])
    for i := last_index; i > 0; i-- {
        sb.WriteString(")")
    }
}

// posts.user.father.mentor.name = "shachar"
// SELECT name FROM (SELECT name FROM mentors WHERE id = (SELECT mentor FROM fathers WHERE id = (SELECT father FROM users WHERE id = posts.user)))

// posts.user.mentor.name = "shachar"
// (SELECT name FROM (SELECT name FROM mentors WHERE id = (SELECT mentor FROM users WHERE id = posts.user))) == "shachar"


// posts.user.name == 'shachar'
// (SELECT name FROM (SELECT name FROM users WHERE id = posts.user)) == 'shachar'
