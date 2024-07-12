package querylang

import (
	"strings"

	querylang_types "github.com/shachar1236/Baasa/query_lang/types"
)

type lexer struct {
	str string
	i   int
}

func (this *lexer) Next() querylang_types.Token {
	start := this.i

	if start >= len(this.str) {
		return querylang_types.Token{
			Type: querylang_types.TOKEN_EOF,
		}
	}

	for is_white_space(this.str[start]) {
		start++
		if start > len(this.str) {
			return querylang_types.Token{
				Type: querylang_types.TOKEN_EOF,
			}
		}
	}

	if this.str[start] == '(' {
		this.i = start + 1
		return querylang_types.Token{
			Value: querylang_types.TokenValueString("("),
			Type:  querylang_types.TOKEN_OPEN_PARENTHESIS,
		}
	}

	if this.str[start] == ')' {
		this.i = start + 1
		return querylang_types.Token{
			Value: querylang_types.TokenValueString(")"),
			Type:  querylang_types.TOKEN_CLOSE_PARENTHESIS,
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
			return querylang_types.Token{
				Value: querylang_types.TokenValueString(my_token),
				Type:  querylang_types.TOKEN_STRING_TYPE,
			}
		}
	}

	// checking if token is operator or collection
	for _, operator := range querylang_types.Filter_lang_operators {
		if my_token == operator {
			return querylang_types.Token{
				Value: querylang_types.TokenValueString(my_token),
				Type:  querylang_types.TOKEN_OPERATOR,
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
		return querylang_types.Token{
			Value: querylang_types.TokenValueString(my_token),
			Type:  querylang_types.TOKEN_NUMBER_TYPE,
		}
	}

	return querylang_types.Token{
		Value: querylang_types.TokenValueVariable{
			Parts: strings.Split(my_token, "."),
		},
		Type: querylang_types.TOKEN_VARIABLE_TYPE,
	}
}
