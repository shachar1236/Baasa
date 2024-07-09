package access_rules

import (
	"context"
	"strings"

	"github.com/shachar1236/Baasa/database/types"
)

// filter language:
// you have the operators: ==, !=, >, >=, <, <=, &&, ||, ()
// special operators:
// '~' - sql like
// '!~' - sql like reverse
// '.' - contains / at least one equel
// '!.' - at least one is not equel
// '>.' - at least one is less
// '>=.' - at least one is less or equel
// '<.' - at least one is greater
// '<=.' - at least one is greater or equel
// '~.' - at least one is sql like
// '!~.' - at least one is sql like reverse

const numbers = "0123456789."

func is_white_space(c byte) bool {
    return c == ' ' || c == '\n'
}

type lexer struct {
    str string
    i int
}

func (this *lexer) Next() types.Token {
    start := this.i

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
            Value: "(",
            Type: types.TOKEN_OPEN_PARENTHESIS,
        }
    }

    if this.str[start] == ')' {
        this.i = start + 1
        return types.Token{
            Value: ")",
            Type: types.TOKEN_CLOSE_PARENTHESIS,
        }
    }

    end := start
    for end < len(this.str) && !is_white_space(this.str[end]) {
        end++
    }

    my_token := this.str[start:end]
    this.i = end

    // checking if token is operator or collection
    for _, operator := range filter_lang_operators {
        if my_token == operator {
            return types.Token{
                Value: my_token,
                Type: types.TOKEN_OPERATOR,
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
            Value: my_token,
            Type: types.TOKEN_NUMBER_TYPE,
        }
    }

    return types.Token{
        Value: my_token,
        Type: types.TOKEN_VARIABLE_TYPE,
    }
}


func (this *AccessRules) AnalyzeUserFilter(filter string) (used_collections []types.Collection, valid bool, tokens []types.Token) {
    valid = true
    lex := lexer{
        str: filter,
        i: 0,
    }

    last_important_token := types.Token{
        Type: types.TOKEN_EOF,
    }

    for  {
        // getting token
        curr_token := lex.Next()

        // check if its two operator in a row or two variables in a row
        if last_important_token .Type != types.TOKEN_EOF {
            if last_important_token.Type == types.TOKEN_OPERATOR && curr_token.Type == types.TOKEN_OPERATOR {
                valid = false
                return
            }
            if last_important_token.Type != types.TOKEN_OPERATOR && curr_token.Type != types.TOKEN_OPERATOR {
                valid = false
                return
            }
        } else {
            if curr_token.Type == types.TOKEN_OPERATOR || curr_token.Type == types.TOKEN_EOF {
                valid = false
                return
            }
        }

        // if last
        if curr_token.Type == types.TOKEN_EOF {
            return
        }

        // checking if its a variable
        if curr_token.Type == types.TOKEN_VARIABLE_TYPE {
            // its a variable
            variable_parts := strings.Split(curr_token.Value, ".")
            if len(variable_parts) > 2 {
                // its nested collections
                for i, collection_name := range variable_parts {
                    if i != len(variable_parts) - 1 {
                        collection, err := this.db.GetCollectionByName(context.Background(), collection_name)
                        if err != nil {
                            valid = false
                            this.logger.Error("Cannot get collection: " + err.Error())
                            return
                        }
                        used_collections = append(used_collections, collection)

                        // checking if collection has the field
                        if !DoesCollectionHasField(collection, variable_parts[i+1]) {
                            valid = false
                            return
                        }
                    }
                }
            } else if len(variable_parts) == 2 {
                // its a field from the collection
            } else {
                valid = false
            }
        }

        if curr_token.Type == types.TOKEN_OPEN_PARENTHESIS || curr_token.Type == types.TOKEN_CLOSE_PARENTHESIS {
            last_important_token = curr_token
        }
    }
}

func DoesCollectionHasField(collection types.Collection, field_name string) bool {
    for _, field := range collection.Fields {
        if field.FieldName == field_name {
            return true
        }
    }
    return false
}
