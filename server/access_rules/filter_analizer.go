package access_rules

import (
	"context"
	"strings"

	"github.com/shachar1236/Baasa/database/types"
)

// filter language:
// you have the operators: =, !=, >, >=, <, <=, &&, ||, ()
// special operators:
// '~' - sql like
// '!~' - sql not like
// '.' - contains / at least one equel
// '!.' - not contains
// '>.' - all are less
// '>=.' - all are less or equel
// '<.' - all are greater
// '<=.' - all are greater or equel
// '!=.' - all are not equel
// '.>' - at least one is less
// '.>=' - at least one is less or equel
// '.<' - at least one is greater
// '.<=' - at least one is greater or equel
// '.!=' - at least one is not equel

const numbers = "0123456789."

func is_white_space(c byte) bool {
	return c == ' ' || c == '\n'
}

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

type lexer struct {
	str string
	i   int
}

func (this *lexer) Next() types.Token {
	start := this.i

	if start >= len(this.str) {
		return types.Token{
			Type: types.TOKEN_EOF,
		}
	}

	for is_white_space(this.str[start]) {
		start++
		if start > len(this.str) {
			return types.Token{
				Type: types.TOKEN_EOF,
			}
		}
	}

	if this.str[start] == '(' {
		this.i = start + 1
		return types.Token{
			Value: types.TokenValueString("("),
			Type:  types.TOKEN_OPEN_PARENTHESIS,
		}
	}

	if this.str[start] == ')' {
		this.i = start + 1
		return types.Token{
			Value: types.TokenValueString(")"),
			Type:  types.TOKEN_CLOSE_PARENTHESIS,
		}
	}

	end := start
	for end < len(this.str) && !is_white_space(this.str[end]) {
		end++
	}

	my_token := this.str[start:end]
	this.i = end

	// checking if token is string
	if len(my_token) > 1 {
		if (my_token[0] == '\'' && my_token[len(my_token)-1] == '\'') || (my_token[0] == '"' && my_token[len(my_token)-1] == '"') {
			return types.Token{
				Value: types.TokenValueString(my_token),
				Type:  types.TOKEN_STRING_TYPE,
			}
		}
	}

	// checking if token is operator or collection
	for _, operator := range types.Filter_lang_operators {
		if my_token == operator {
			return types.Token{
				Value: types.TokenValueString(my_token),
				Type:  types.TOKEN_OPERATOR,
			}
		}
	}

	// checking if token is a number
	is_number := true
	for _, c := range my_token {
		is_inside := false
		for _, number := range numbers {
			if c == number {
				is_inside = true
				break
			}
		}

		if !is_inside {
			is_number = false
			break
		}
	}

	if is_number {
		return types.Token{
			Value: types.TokenValueString(my_token),
			Type:  types.TOKEN_NUMBER_TYPE,
		}
	}

	return types.Token{
		Value: types.TokenValueVariable{
			Parts: strings.Split(my_token, "."),
		},
		Type: types.TOKEN_VARIABLE_TYPE,
	}
}

func hasCollection(collections []types.Collection, collection_name string) bool {
	for _, curr_collection := range collections {
		if curr_collection.Name == collection_name {
			return true
		}
	}
	return false
}

func (this *AccessRules) analyzeVariableParts(my_collection_name string, token_as_variable *types.TokenValueVariable) (valid bool, used_collections CollectionsSet) {
	variable_parts := token_as_variable.Parts
	used_collections.Init()
	// its nested collections
	my_collection, err := this.db.GetCollectionByName(context.Background(), variable_parts[0])
	if err != nil {
		valid = false
		this.logger.Error("Cannot get collection: " + err.Error())
		return
	}

	if my_collection.Name != my_collection_name {
		valid = false
		return
	}

	token_as_variable.Fields = make([]types.TokenValueVariablePart, len(variable_parts)-1)

	last_collection := my_collection
	for i := 1; i < len(variable_parts)-1; i++ {
		my_field := types.TableField{
			ID: -1,
		}
		for _, field := range last_collection.Fields {
			if field.FieldName == variable_parts[i] {
				my_field = field
				break
			}
		}

		if my_field.ID != -1 {

            token_as_variable.Fields[i-1].PartType = types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_FIELD_TYPE
            token_as_variable.Fields[i-1].Field = types.TokenValueVariablePartField{
                FieldName: my_field.FieldName,
                FieldCollection: last_collection.Name,
            }

			if my_field.IsForeignKey {
				curr_collection_name := my_field.FkRefersToTable.String
				curr_collection, err := this.db.GetCollectionByName(context.Background(), curr_collection_name)
				if err != nil {
					valid = false
					this.logger.Error("Cannot get collection: " + err.Error())
					return
				}

                token_as_variable.Fields[i-1].Field.FkRefersToCollection = curr_collection_name

				last_collection = curr_collection
				used_collections.Add(curr_collection)
			} else {
				valid = false
				return
			}
		} else {
			// maybe its a list
			if i == len(variable_parts)-2 {
				curr_collection, err := this.db.GetCollectionByName(context.Background(), variable_parts[i])
				if err != nil {
					valid = false
					this.logger.Error("Cannot get collection: " + err.Error())
					return
				}

				// checking if collection has relation to last collection
				has_relation := false
				var relation types.TableField
				for _, field := range curr_collection.Fields {
					if field.IsForeignKey && field.FkRefersToTable.String == last_collection.Name {
						has_relation = true
						relation = field
						break
					}
				}

				if !has_relation {
					valid = false
					return
				}

                token_as_variable.Fields[i-1].PartType = types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_TYPE
                token_as_variable.Fields[i-1].Collection = struct{CollectionName string;FkToLastPartName string}{
                    CollectionName: curr_collection.Name,
                    FkToLastPartName: relation.FieldName,
                }


				last_collection = curr_collection
				used_collections.Add(curr_collection)
			} else {
				// its not a list
				valid = false
				return
			}
		}
	}

	if !DoesCollectionHasField(last_collection, variable_parts[len(variable_parts)-1]) {
		valid = false
		return
	}

    last_index := len(token_as_variable.Fields) - 1
    token_as_variable.Fields[last_index].PartType = types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_FIELD_TYPE
    token_as_variable.Fields[last_index].Field = types.TokenValueVariablePartField{
        FieldName: variable_parts[len(variable_parts) - 1],
        FieldCollection: last_collection.Name,
    }


	valid = true
	return
}

func (this *AccessRules) AnalyzeUserFilter(my_collection_name string, filter string) (used_collections CollectionsSet, valid bool, tokens []types.Token) {
	valid = true
	lex := lexer{
		str: filter,
		i:   0,
	}

	used_collections.Init()
	tokens = make([]types.Token, 0)

	last_important_token := types.Token{
		Type: types.TOKEN_EOF,
	}

	is_left_side := true

	for {
		// getting token
		curr_token := lex.Next()

		if curr_token.Type == types.TOKEN_EOF {
			break
		}

		// check if its two operator in a row or two variables in a row
		if last_important_token.Type != types.TOKEN_EOF {
			if last_important_token.Type == types.TOKEN_OPERATOR && curr_token.Type == types.TOKEN_OPERATOR {
				valid = false
				return
			}
			if last_important_token.Type != types.TOKEN_OPERATOR && curr_token.Type != types.TOKEN_OPERATOR {
				valid = false
				return
			}
		} else {
			// checking if start token is operator or EOF
			if curr_token.Type == types.TOKEN_OPERATOR || curr_token.Type == types.TOKEN_EOF {
				valid = false
				return
			}
		}

		// if last
		if curr_token.Type == types.TOKEN_EOF {
			return
		}

		// changing side
		if curr_token.Type == types.TOKEN_VARIABLE_TYPE || curr_token.Type == types.TOKEN_NUMBER_TYPE {
			is_left_side = !is_left_side
		}

		// checking if its a variable
		if curr_token.Type == types.TOKEN_VARIABLE_TYPE {
			// its a variable
			token_as_variable := curr_token.Value.(types.TokenValueVariable)
			if len(token_as_variable.Parts) > 2 {
				my_valid, added_used_collections := this.analyzeVariableParts(my_collection_name, &token_as_variable)
				if !my_valid {
					valid = false
					return
				}
				curr_token.Value = token_as_variable
				used_collections.Union(added_used_collections)
			} else if len(token_as_variable.Parts) == 2 {
				// its a field from the collection
				if token_as_variable.Parts[0] == my_collection_name {
					// its nested collections
					my_collection, err := this.db.GetCollectionByName(context.Background(), token_as_variable.Parts[0])
					if err != nil {
						valid = false
						this.logger.Error("Cannot get collection: " + err.Error())
						return
					}

					if !DoesCollectionHasField(my_collection, token_as_variable.Parts[1]) {
						valid = false
						return
					}
				} else {
					valid = false
					return
				}
			} else {
				valid = false
				return
			}
		}

		if curr_token.Type != types.TOKEN_OPEN_PARENTHESIS && curr_token.Type != types.TOKEN_CLOSE_PARENTHESIS {
			last_important_token = curr_token
		}

		tokens = append(tokens, curr_token)
	}

	return
}

func DoesCollectionHasField(collection types.Collection, field_name string) bool {
	for _, field := range collection.Fields {
		if field.FieldName == field_name {
			return true
		}
	}
	return false
}
