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
	variable_parts := token_as_variable.Parts
	last_index := len(variable_parts) - 1
	// sb.WriteString("(SELECT ")
	// sb.WriteString(variable_parts[last_index])
	// sb.WriteString(" FROM ")
	for i := last_index; i > 1; i-- {
		sb.WriteString(" (SELECT ")
        field_name := variable_parts[i]
        if token_as_variable.ExtandsToList && i == last_index - 1 {
            field_name = token_as_variable.PartToCollectionField[variable_parts[i-1]].FkFieldName.String
        }
		sb.WriteString(field_name)
		sb.WriteString(" FROM ")
		curr_used_collection := token_as_variable.PartToCollection[variable_parts[i-1]].Name
		sb.WriteString(curr_used_collection)
		sb.WriteString(" WHERE ")
		filters := used_collections_filters[curr_used_collection]
		if filters != "" {
			sb.WriteString("(")
			sb.WriteString(filters)
			sb.WriteString(") AND ")
		}
        fk_variable := token_as_variable.PartToCollectionField[variable_parts[i-1]].FkFieldName.String
        if token_as_variable.ExtandsToList && i == last_index {
            fk_variable = token_as_variable.PartToCollectionField[variable_parts[i-1]].FieldName
        }
		sb.WriteString(fk_variable)
		sb.WriteString(" = ")
	}
	sb.WriteString(variable_parts[0])
	sb.WriteString(".")
	sb.WriteString(variable_parts[1])
	for i := last_index; i > 1; i-- {
        sb.WriteString(" LIMIT 1)")
	}
}

// posts.user.father.mentor.name = "shachar"
// SELECT name FROM (SELECT name FROM mentors WHERE id = (SELECT mentor FROM fathers WHERE id = (SELECT father FROM users WHERE id = posts.user)))

// posts.user.mentor.name = "shachar"
// (SELECT name FROM mentors WHERE id = (SELECT mentor FROM users WHERE id = posts.user)) == "shachar"

// --->
// posts.user.name == 'shachar'
// (SELECT name FROM (SELECT name FROM users WHERE id = posts.user)) == 'shachar'
// (SELECT name FROM users WHERE id = posts.user) == 'shachar'

// 'test' . Posts.user.title.Posts - no!!!!!!!!!!!!!! 
// 'test' . Posts.user.Posts.title
// 'test' IN (SELECT title FROM Posts WHERE user = (SELECT id FROM Users WHERE id = Posts.user))

// <------
// 'test' . Posts.user.Comments.content
// 'test' IN (SELECT content FROM Comments WHERE created_by_user = (SELECT id FROM users WHERE name = Posts.user))

// 'test' . Posts.user.mentor.Posts.title
// 'test' IN (SELECT title FROM Posts WHERE user = (SELECT mentor FROM Users WHERE id = (SELECT id FROM users WHERE name = Posts.user))
