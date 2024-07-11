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

func (db *SqliteDB) BuildUserCustomQuery(
	collection_name string,
	fields []string,
	filter_tokens []types.Token,
	sort_by string,
	expand string,
	used_collections_filters map[string]string,
) (query string, err error) {
	_ = sqlf.Select(strings.Join(fields, ",")).From(collection_name)

	sb_mutex.Lock()
	defer sb_mutex.Unlock()
	defer sb.Reset()

	for _, token := range filter_tokens {
		if token.Type == types.TOKEN_OPERATOR {
			// TODO: change for supporting special operators
			token_as_string := token.Value.(types.TokenValueString)
			// sb.WriteString(string(token_as_string))
			switch token_as_string {
			case "&&":
				sb.WriteString(" AND ")
				break
			case "||":
				sb.WriteString(" OR ")
				break
			case "~":
				sb.WriteString(" LIKE ")
				break
			case "!~":
				sb.WriteString(" NOT LIKE ")
				break
			case ".":
				sb.WriteString(" IN ")
				break
			case "!.":
				sb.WriteString(" NOT IN ")
				break
			case ">.":
				sb.WriteString(" > ALL ")
				break
			case ">=.":
				sb.WriteString(" >= ALL ")
				break
			case "<.":
				sb.WriteString(" < ALL ")
				break
			case "<=.":
				sb.WriteString(" <= ALL ")
				break
			case "!=.":
				sb.WriteString(" != ALL ")
				break
			case ".>":
				sb.WriteString(" > ANY ")
				break
			case ".>=":
				sb.WriteString(" >= ANY ")
				break
			case ".<":
				sb.WriteString(" < ANY ")
				break
			case ".<=":
				sb.WriteString(" <= ANY ")
				break
			case ".!=":
				sb.WriteString(" != ANY ")
				break
			default:
				sb.WriteString(" ")
				sb.WriteString(string(token_as_string))
				sb.WriteString(" ")
				break
			}
		} else if token.Type == types.TOKEN_VARIABLE_TYPE {
			token_as_variable := token.Value.(types.TokenValueVariable)
			variable_parts := token_as_variable.Parts
			if len(variable_parts) > 2 {
				// its nested collections
				createExpandedSelect(&token_as_variable, &sb, used_collections_filters)
			} else if len(variable_parts) == 2 {
				// its a field from the collection
				sb.WriteString(token_as_variable.Parts[0])
				sb.WriteString(".")
				sb.WriteString(token_as_variable.Parts[1])
			} else {
				// not valid
			}
		} else {
			token_as_string := token.Value.(types.TokenValueString)
			sb.WriteString(string(token_as_string))
		}
		sb.WriteString(" ")
	}

	query = sb.String()

	return
}

// build and runs user query
func (db *SqliteDB) RunUserCustomQuery(
	collection_name string,
	fields []string,
	filter_tokens []types.Token,
	sort_by string,
	expand string,
	used_collections_filters map[string]string, // map[collection_name]filters
) (resJson string, err error) {
	query, err := db.BuildUserCustomQuery(collection_name, fields, filter_tokens, sort_by, expand, used_collections_filters)
	if err != nil {
		db.logger.Error("Error in user custom query: " + err.Error())
		return
	}
	fmt.Println("Printing user custom query")
	fmt.Println()
	fmt.Println()
	fmt.Println(query)
	fmt.Println()
	fmt.Println()
	return
}

func createExpandedSelect(token_as_variable *types.TokenValueVariable, sb *strings.Builder, used_collections_filters map[string]string) {
	last_index := len(token_as_variable.Fields) - 1

	curr_field := token_as_variable.Fields[last_index]
	sb.WriteString("(SELECT ")
	sb.WriteString(curr_field.FieldName)
	sb.WriteString(" FROM ")
	sb.WriteString(curr_field.FieldTable)
	sb.WriteString(" WHERE ")
	for i := last_index - 1; i > 0; i-- {
		curr_field = token_as_variable.Fields[i]
		field_name := curr_field.FieldName
		field_table := curr_field.FieldTable
		relation_field := "id"
		// if curr_field.ViaChildToParent {
			// relation_field = field_name
			// field_name = "id"
		// }
		sb.WriteString(relation_field)
		sb.WriteString(" = (SELECT ")
		sb.WriteString(field_name)
		sb.WriteString(" FROM ")
		sb.WriteString(field_table)
		sb.WriteString(" WHERE ")
	}

	curr_field = token_as_variable.Fields[0]
	sb.WriteString(token_as_variable.LastFieldRefersToField)
	sb.WriteString(" = ")
	sb.WriteString(curr_field.FieldName)

	for i := 0; i < last_index; i++ {
		curr_field = token_as_variable.Fields[i]
		if !curr_field.ViaChildToParent {
			sb.WriteString(" LIMIT 1)")
		} else {
			sb.WriteString(")")
		}
	}

}

// posts.user.father.mentor.name = "shachar"
// (SELECT name FROM mentors WHERE id = (SELECT mentor FROM fathers WHERE id = (SELECT father FROM users WHERE id = posts.user))
// { "field_name": "posts.user", "table": "posts",  "field_foreign_key_to": ""}

// { "field_name": "father", "table": "users",  "field_foreign_key_to": ["fathers", "id"]}
// { "field_name": "mentor", "table": "fathers",  "field_foreign_key_to": ["mentors", "id"]}

// { "field_name": "name", "table": "mentors",  "field_foreign_key_to": ""}

// posts.user.mentor.name = "shachar"
// (SELECT name FROM mentors WHERE id = (SELECT mentor FROM users WHERE id = posts.user)) == "shachar"
// fields =

// --->
// posts.user.name == 'shachar'
// (SELECT name FROM users WHERE id = posts.user) == 'shachar'
// fields =
// { "field_name": "Posts.user", "table": "Posts",  "field_foreign_key_to": []}
// { "field_name": "name", "table": "users",  "field_foreign_key_to": []}

// 'test' . Posts.user.Comments.content
// 'test' IN (SELECT content FROM Comments WHERE created_by_user = (SELECT id FROM users WHERE name = Posts.user))
// fields =
// { "field_name": "Posts.user", "table": "Posts",  "field_foreign_key_to": []}

// { "field_name": "id", "table": "users",  "field_foreign_key_refers_to": ["Comments", "created_by_user"]} ||
//      { "field_name": "created_by_user", "table": "Comments", via_child_to_parent = true,  "field_foreign_key_refers_to": ["Users", "id"]}

// { "field_name": "content", "table": "Comments",  "field_foreign_key_to": ""}

// ------------
// { "field_name": "Posts.user", "table": "Posts",  "field_foreign_key_to": []}
// { "field_name": "content", "table": "Comments",  "field_foreign_key_to": ""}

// 'test' . Posts.user.mentor.Posts.title
// 'test' IN (SELECT title FROM Posts WHERE user = (SELECT mentor FROM Users WHERE id = Posts.user))
// { "field_name": "Posts.user", "table": "Posts",  "field_foreign_key_to": ""}
// { "field_name": "mentor", "table": "Users",  "field_foreign_key_to": "mentors"}
// { "field_name": "title", "table": "Posts",  "field_foreign_key_to": ""}


// Posts.user
// { "field_name": "mentor", "table": "Users",  "field_foreign_key_to": "id"}
// { "field_name": "title", "table": "Posts",  "field_foreign_key_to": "user"}
