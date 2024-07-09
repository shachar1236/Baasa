package sqlite

import (
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
    used_collections []types.Collection,
    used_collections_filters []string,
) (resJson string, err error) {
    _ = sqlf.Select(strings.Join(fields, ",")).From(collection_name)

    sb_mutex.Lock()
    defer sb_mutex.Unlock()
    defer sb.Reset()

    for _, token := range filter_tokens {
        if token.Type == types.TOKEN_OPERATOR {

        } else if token.Type == types.TOKEN_VARIABLE_TYPE {
            variable_parts := strings.Split(token.Value, ".")
            if len(variable_parts) > 2 {
                // its nested collections
                for i, collection_name := range variable_parts {
                    if i != len(variable_parts) - 1 {
                    }
                }
            } else if len(variable_parts) == 2 {
                // its a field from the collection
                sb.WriteString(token.Value)
            } else {
                // not valid
            }
        } else {
            sb.WriteString(token.Value)
        }
        sb.WriteString(" ")
    }
}

func createExpandedSelect(variable_parts []string, sb strings.Builder) {
    // TODO: add collection filters
    last_index := len(variable_parts) - 1
    sb.WriteString("SELECT ")
    sb.WriteString(variable_parts[last_index])
    sb.WriteString(" FROM (SELECT ")
    sb.WriteString(variable_parts[last_index])
    sb.WriteString(" FROM ")
}

// posts.user.father.mentor.name = "shachar"
// SELECT name FROM (SELECT name FROM mentors WHERE id = (SELECT mentor FROM fathers WHERE id = (SELECT father FROM users WHERE id = posts.user)))

// posts.user.mentor.name = "shachar"
// (SELECT name FROM (SELECT name FROM mentors WHERE id = (SELECT mentor FROM users WHERE id = posts.user))) == "shachar"


// posts.user.name == 'shachar'
// (SELECT name FROM (SELECT name FROM users WHERE id = posts.user)) == 'shachar'
