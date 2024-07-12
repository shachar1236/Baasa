package sqlite

import (
	"fmt"
	"strings"
	"sync"

	"github.com/leporo/sqlf"
    "github.com/deckarep/golang-set/v2"
	querylang_types "github.com/shachar1236/Baasa/query_lang/types"
)

var sb = strings.Builder{}
var sb_mutex = sync.Mutex{}

func (db *SqliteDB) BuildUserCustomQuery(
	collection_name string,
	fields []string,
	filter_tokens []querylang_types.Token,
	sort_by string,
    analyzed_expand []querylang_types.TokenValueVariable,
	used_collections_filters map[string]string,
) (where_query string, err error) {

	sb_mutex.Lock()
	defer sb_mutex.Unlock()
	defer sb.Reset()

	for _, token := range filter_tokens {
		if token.Type == querylang_types.TOKEN_OPERATOR {
			// TODO: change for supporting special operators
			token_as_string := token.Value.(querylang_types.TokenValueString)
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
		} else if token.Type == querylang_types.TOKEN_VARIABLE_TYPE {
			token_as_variable := token.Value.(querylang_types.TokenValueVariable)
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
			token_as_string := token.Value.(querylang_types.TokenValueString)
			sb.WriteString(string(token_as_string))
		}
		sb.WriteString(" ")
	}

	where_query = sb.String()
    sb.Reset()

    select_fields, expand_filters, err := joinExpandedFields(collection_name, analyzed_expand, fields, &sb, used_collections_filters)
    if err != nil {
        return "", err
    }
    join_query := sb.String()
    sb.Reset()

    sql_query := sqlf.Select(strings.Join(fields, ",")).Select(strings.Join(select_fields, ",")).From(collection_name)
    sql_query.Expr(join_query)

    // adding where
    sql_query = sql_query.Where(where_query)
    for _, filter := range expand_filters {
        sql_query = sql_query.Where(filter)
    }
    filters := used_collections_filters[collection_name]
    if filters != "" {
        sql_query = sql_query.Where(filters)
    }

	return sql_query.String(), nil
}

// build and runs user query
func (db *SqliteDB) RunUserCustomQuery(
	collection_name string,
	fields []string,
	filter_tokens []querylang_types.Token,
	sort_by string,
	analyzed_expand []querylang_types.TokenValueVariable,
	used_collections_filters map[string]string, // map[collection_name]filters
) (resJson string, err error) {
	query, err := db.BuildUserCustomQuery(collection_name, fields, filter_tokens, sort_by, analyzed_expand, used_collections_filters)
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

func createExpandedSelect(token_as_variable *querylang_types.TokenValueVariable, sb *strings.Builder, used_collections_filters map[string]string) {
	last_index := len(token_as_variable.Fields) - 1
	curr_field := token_as_variable.Fields[last_index].Field
	sb.WriteString("(SELECT ")
	sb.WriteString(curr_field.FieldName)
	sb.WriteString(" FROM ")
	sb.WriteString(curr_field.FieldCollection)
	sb.WriteString(" WHERE ")
	filters := used_collections_filters[curr_field.FieldCollection]
	if filters != "" {
		sb.WriteString(filters)
		sb.WriteString(" AND ")
	}

	startIndex := last_index - 1
	curr_part := token_as_variable.Fields[last_index-1]
	closingParenthesseis := " LIMIT 1)"
	if curr_part.PartType == querylang_types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_TYPE {
		sb.WriteString(curr_part.Collection.FkToLastPartName)
		sb.WriteString(" = ")
		startIndex--
		closingParenthesseis = ")"
	} else {
		sb.WriteString("id = ")
	}

	for i := startIndex; i > 0; i-- {
		curr_field := token_as_variable.Fields[i].Field
		sb.WriteString("(SELECT ")
		sb.WriteString(curr_field.FieldName)
		sb.WriteString(" FROM ")
		sb.WriteString(curr_field.FieldCollection)
		sb.WriteString(" WHERE ")
		filters := used_collections_filters[curr_field.FieldCollection]
		if filters != "" {
			sb.WriteString(filters)
			sb.WriteString(" AND ")
		}
		sb.WriteString("id =  ")
	}

	curr_part = token_as_variable.Fields[0]
	curr_field = curr_part.Field
	filters = used_collections_filters[curr_field.FkRefersToCollection]
	has_filters := filters != ""
	if has_filters {
		sb.WriteString("(SELECT id FROM ")
		sb.WriteString(curr_field.FieldCollection)
		sb.WriteString(" WHERE ")
		sb.WriteString(filters)
		sb.WriteString(" AND ")
		sb.WriteString("id =  ")
		startIndex++
	}
	sb.WriteString(curr_part.Field.FieldCollection)
	sb.WriteString(".")
	sb.WriteString(curr_part.Field.FieldName)

	for i := startIndex; i > 0; i-- {
		sb.WriteString(" LIMIT 1)")
	}

	sb.WriteString(closingParenthesseis)
}

func joinExpandedFields(my_collection_name string, analyzed_expand []querylang_types.TokenValueVariable, fields []string, sb *strings.Builder, used_collections_filters map[string]string) (select_fields []string, where []string, err error) {
    used_collections := mapset.NewSet[string]()
    for _, exp := range analyzed_expand {
        if len(exp.Fields) > 0 {
            last_part := exp.Fields[len(exp.Fields) - 1]
            if !used_collections.Contains(last_part.Field.FieldCollection) {
                sb.WriteString(" LEFT JOIN ")
                sb.WriteString(last_part.Field.FieldCollection)
                sb.WriteString(" ON ")
                sb.WriteString(last_part.Field.FieldCollection)
                sb.WriteString(".id")
                sb.WriteString(" = ")
                var select_token querylang_types.TokenValueVariable
                if len(exp.Fields) > 2 {
                    select_token.Parts = exp.Parts[:len(exp.Parts) - 1]
                    select_token.Fields = exp.Fields[:len(exp.Fields) - 1]
                    var my_sb strings.Builder
                    createExpandedSelect(&select_token, &my_sb, used_collections_filters)
                    sb.WriteString(my_sb.String())
                } else {
                    sb.WriteString(exp.Fields[0].Field.FieldName)
                }

                used_collections.Add(last_part.Field.FieldCollection)
            }

            if last_part.Field.FieldName != "*" {
                select_fields = append(select_fields, last_part.Field.FieldCollection + "." + last_part.Field.FieldName)
            } else {
                
            }
        }
    }

    iter := used_collections.Iter()
    for collection := range iter {
        filters := used_collections_filters[collection]
        if filters != "" {
            where = append(where, filters)
        }
    }

    return
}

// collection = Posts
// fields = [title]
// expanded = [Posts.user.name]

// SELECT Posts.title, Users.Name FROM Posts 
//  LEFT JOIN Users
//  ON Users.id = (SELECT id FROM Users WHERE id = Posts.user AND ...)
