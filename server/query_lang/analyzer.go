package querylang

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/shachar1236/Baasa/database"
	"github.com/shachar1236/Baasa/database/types"
	querylang_types "github.com/shachar1236/Baasa/query_lang/types"
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
// TODO: add ~. and .~

const numbers = "0123456789."

type Analyzer struct {
	logger *slog.Logger
	db     database.Database
}

func New(db database.Database) Analyzer {
	logFile, err := os.OpenFile("logs/query_lang_analyzer.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	var mw io.Writer
	if err != nil {
		mw = os.Stdout
	} else {
		mw = io.MultiWriter(os.Stdout, logFile)
	}
	logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))

	return Analyzer{
		logger: logger,
		db:     db,
	}
}

func (this *Analyzer) AnalyzeVariableParts(my_collection_name string, token_as_variable *querylang_types.TokenValueVariable, analyze_options int) (valid bool, used_collections querylang_types.CollectionsSet) {
	variable_parts := token_as_variable.Parts
	used_collections.Init()
	if analyze_options&querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_OPTIONS_VALIDATE_VARIABLE != 0 {
		if len(variable_parts) == 0 {
			return false, used_collections
		}
		// for _, part := range variable_parts {
		// TODO: validate variable if contains valid charecters without whitespace, for more info check https://trello.com/c/KYiiuBOM
		// }
	}

	var err error

	// is_analyzing_filter := analyze_type == querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_TYPE_FILTER
	is_analyzing_join := analyze_options&querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_OPTIONS_JOIN != 0
	analize_options_not_list := analyze_options&querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_OPTIONS_NOT_LIST != 0
	analize_options_dont_fix := analyze_options&querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_OPTIONS_DONT_FIX != 0
	backward_expand_occured := false

	var my_collection types.Collection
	if variable_parts[0] == my_collection_name {
		my_collection, err = this.db.GetCollectionByName(context.Background(), variable_parts[0])
		if err != nil {
			valid = false
			this.logger.Error("Cannot get collection: " + err.Error())
			return
		}
	} else if !analize_options_dont_fix {
		my_collection, err = this.db.GetCollectionByName(context.Background(), my_collection_name)
		if err != nil {
			valid = false
			this.logger.Error("Cannot get collection: " + err.Error())
			return
		}
		parts := make([]string, len(variable_parts)+1)
		parts[0] = my_collection_name
		for i := 0; i < len(variable_parts); i++ {
			parts[1+i] = variable_parts[i]
		}
		variable_parts = parts
		token_as_variable.Parts = parts
	} else {
		valid = false
		return
	}

	token_as_variable.Fields = make([]querylang_types.TokenValueVariablePart, len(variable_parts)-1)

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

			token_as_variable.Fields[i-1].PartType = querylang_types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_FIELD_TYPE
			token_as_variable.Fields[i-1].Field = querylang_types.TokenValueVariablePartField{
				FieldName:       my_field.FieldName,
				FieldCollection: last_collection.Name,
			}

			if my_field.IsForeignKey {
				curr_collection_name := my_field.FkRefersToTable.String
				curr_collection, exists := used_collections.GetCollectionByName(curr_collection_name)
				if !exists {
					curr_collection, err = this.db.GetCollectionByName(context.Background(), curr_collection_name)
					if err != nil {
						valid = false
						this.logger.Error("Cannot get collection: " + err.Error())
						return
					}
				}

				token_as_variable.Fields[i-1].Field.FkRefersToCollection = curr_collection_name
				token_as_variable.Fields[i-1].Field.FieldCollectionPointer = &curr_collection

				last_collection = curr_collection
				used_collections.Add(curr_collection)
			} else {
				valid = false
				return
			}
		} else if !analize_options_not_list {
			// maybe its a list
			if i == len(variable_parts)-2 && !backward_expand_occured {
				curr_collection, exists := used_collections.GetCollectionByName(variable_parts[i])
				if !exists {
					curr_collection, err = this.db.GetCollectionByName(context.Background(), variable_parts[i])
					if err != nil {
						valid = false
						this.logger.Error("Cannot get collection: " + err.Error())
						return
					}
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

				token_as_variable.Fields[i-1].PartType = querylang_types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_TYPE
				token_as_variable.Fields[i-1].Collection = struct {
					CollectionName    string
					FkToLastPartField types.TableField
				}{
					CollectionName:    curr_collection.Name,
					FkToLastPartField: relation,
				}

				last_collection = curr_collection
				used_collections.Add(curr_collection)
				backward_expand_occured = true
			} else {
				// its not a list
				valid = false
				return
			}
		}
	}

	if variable_parts[len(variable_parts)-1] == "*" {
		if !is_analyzing_join {
			valid = false
			return
		}
	} else if !DoesCollectionHasField(last_collection, variable_parts[len(variable_parts)-1]) {
		valid = false
		return
	}

	last_index := len(token_as_variable.Fields) - 1
	token_as_variable.Fields[last_index].PartType = querylang_types.TOKEN_VALUE_VARIABLE_PART_COLLECTION_FIELD_TYPE
	token_as_variable.Fields[last_index].Field = querylang_types.TokenValueVariablePartField{
		FieldName:       variable_parts[len(variable_parts)-1],
		FieldCollection: last_collection.Name,
	}
	token_as_variable.Fields[last_index].Field.FieldCollectionPointer = &last_collection

	valid = true
	return
}

func (this *Analyzer) AnalyzeUserFilter(my_collection_name string, filter string) (used_collections querylang_types.CollectionsSet, valid bool, tokens []querylang_types.Token) {
	valid = true
	lex := lexer{
		str: filter,
		i:   0,
	}

	used_collections.Init()
	tokens = make([]querylang_types.Token, 0)

	last_important_token := querylang_types.Token{
		Type: querylang_types.TOKEN_EOF,
	}

	is_left_side := true

	for {
		// getting token
		curr_token := lex.Next()

		if curr_token.Type == querylang_types.TOKEN_EOF {
			break
		}

		// check if its two operator in a row or two variables in a row
		if last_important_token.Type != querylang_types.TOKEN_EOF {
			if last_important_token.Type == querylang_types.TOKEN_OPERATOR && curr_token.Type == querylang_types.TOKEN_OPERATOR {
				valid = false
				return
			}
			if last_important_token.Type != querylang_types.TOKEN_OPERATOR && curr_token.Type != querylang_types.TOKEN_OPERATOR {
				valid = false
				return
			}
		} else {
			// checking if start token is operator or EOF
			if curr_token.Type == querylang_types.TOKEN_OPERATOR || curr_token.Type == querylang_types.TOKEN_EOF {
				valid = false
				return
			}
		}

		// if last
		if curr_token.Type == querylang_types.TOKEN_EOF {
			return
		}

		// changing side
		if curr_token.Type == querylang_types.TOKEN_VARIABLE_TYPE || curr_token.Type == querylang_types.TOKEN_NUMBER_TYPE {
			is_left_side = !is_left_side
		}

		// checking if its a variable
		if curr_token.Type == querylang_types.TOKEN_VARIABLE_TYPE {
			// its a variable
			token_as_variable := curr_token.Value.(querylang_types.TokenValueVariable)
			my_valid, added_used_collections := this.AnalyzeVariableParts(my_collection_name, &token_as_variable, querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_OPTIONS_FILTER)
			if !my_valid {
				this.logger.Warn("Invalid token variable")
				valid = false
				return
			}
			curr_token.Value = token_as_variable
			used_collections.Union(added_used_collections)
		}

		if curr_token.Type != querylang_types.TOKEN_OPEN_PARENTHESIS && curr_token.Type != querylang_types.TOKEN_CLOSE_PARENTHESIS {
			last_important_token = curr_token
		}

		tokens = append(tokens, curr_token)
	}

	return
}
