package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
	"github.com/shachar1236/Baasa/database/types"
	querylang "github.com/shachar1236/Baasa/query_lang"
	querylang_types "github.com/shachar1236/Baasa/query_lang/types"
)

const SEARCH_RULES_FILENAME = "search.lua"
const BASE_LIMIT = 1000
const MAX_LIMIT = 10000

func addRoutes(mux *http.ServeMux, logger *slog.Logger, db database.Database, access_rules *access_rules.AccessRules, query_lang_analyzer *querylang.Analyzer) {
	mux.Handle("/query", handleQuery(logger, db, access_rules))
	mux.Handle("/collection/search", handleCollectionSearch(logger, db, access_rules, query_lang_analyzer))
	mux.Handle("/", http.NotFoundHandler())
}

func createRequestObject(r *http.Request, user types.User) (access_rules.Request, error) {
    if user.ID == -1 {
        return access_rules.Request{}, nil
    }
    request := access_rules.Request{
        Method:  r.Method,
        Auth:    user,
    }
    for k, v := range r.Header {
        if len(v) > 0 {
            request.Headers[k] = v[0]
        }
    }
    err := r.ParseForm()
    if err != nil {
        return request, err
    }
    for k, v := range r.Form {
        if len(v) > 0 {
            request.Data[k] = v[0]
        }
    }

    return request, err
}

func handleQuery(logger *slog.Logger, db database.Database, ar *access_rules.AccessRules) http.Handler {

    type QueryMessage struct {
        QueryId int64 `json:"QueryId"`
        Session string `json:"Session"`
        QueryArgs map[string]any `json:"QuaryArgs"`
    }

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            msg, err := decode[QueryMessage](r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			query_filters := ""

			// getting query
			query, err := db.GetQuaryById(r.Context(), msg.QueryId)
			if err != nil {
				logger.Error("Cannot get query: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// getting user
			user, err := db.GetUserBySession(r.Context(), msg.Session)
			if err != nil {
				logger.Error("Cannot get user")
				user.ID = -1
				user.Username = ""
			}

			// check if query is by the rules
            request, err := createRequestObject(r, user)
            if err != nil {
                logger.Error("Cannot create request object: " + err.Error())
                return
            }

			accept, err := ar.CheckRules(query.QueryRulesFilePath, &query_filters, request, msg.QueryArgs)
			if err != nil {
				logger.Error("Cannot check query rules: ", err)
				http.Error(w, "An error occured", http.StatusInternalServerError)
				return
			}

			if accept {
				// run query
				resJson, err := db.RunQueryWithFilters(r.Context(), query, msg.QueryArgs, query_filters)
				if err != nil {
					logger.Error("Cannot run query: ", err)
					http.Error(w, "An error occured", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(resJson)
			}
		},
	)
}


func handleCollectionSearch(logger *slog.Logger, db database.Database, ar *access_rules.AccessRules, query_lang_anayzer *querylang.Analyzer) http.Handler {
    // TODO: return better errors to user

    // TODO: validate each field in the message to see if there is an sql injection in it.
    type SearchMessage struct {
        CollectionName string `json:"CollectionName"`
        Session string `json:"Session"`
        Fields []string `json:"Fields"`
        Filter string `json:"Filter"`
        SortBy []string `json:"SortBy"`
        Expand []string `json:"Expand"`
        Limit int `json:"Limit"`
        Offset int `json:"Offset"`
    }

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            var err error
            msg, err := decode[SearchMessage](r)
			if err != nil {
                logger.Error("Got error while trying to decode SearchMessage: " + err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

            // getting collection
            collection, err := db.GetCollectionByName(r.Context(), msg.CollectionName)
            if err != nil {
                logger.Error("Cannot get collection: " + err.Error())
                http.Error(w, "", http.StatusInternalServerError)
                return
            }

			// getting user
			user, err := db.GetUserBySession(r.Context(), msg.Session)
			if err != nil {
				logger.Error("Cannot get user")
				user.ID = -1
				user.Username = ""
			}

            // checking if fields exists
            for _, field_name := range msg.Fields {
                exists := false
                if field_name == "id" {
                    exists = true
                } else {
                    for _, field2 := range collection.Fields {
                        if field_name == field2.FieldName {
                            exists = true
                            break
                        }
                    }
                }

                if !exists {
                    logger.Error("Field does not exist: " + field_name)
                    w.WriteHeader(http.StatusBadRequest)
                    return
                }
            }

            // checks if SortBy valid
            for _, field_name := range msg.SortBy {
                exists := false
                if field_name[0] == '+' || field_name[0] == '-' {
                    field_name = field_name[1:]
                }
                if field_name == "id" {
                    exists = true
                } else {
                    for _, field2 := range collection.Fields {
                        if field_name == field2.FieldName {
                            exists = true
                        }
                    }
                }

                if !exists {
                    logger.Error("Field does not exist: " + field_name)
                    w.WriteHeader(http.StatusBadRequest)
                    return
                }
            }

            // checks if limit is valid
            msg_limit := msg.Limit

            if msg_limit < 1 {
                msg_limit = BASE_LIMIT
            }

            if msg_limit > MAX_LIMIT {
                msg_limit = MAX_LIMIT
            }

            if msg_limit == 0 {
                msg_limit = BASE_LIMIT
            }

            // checks if offset is valid
            msg_offset := msg.Offset

            if msg_offset >= msg_limit {
                w.WriteHeader(http.StatusBadRequest)
                return
            }

			// check if query is by the rules

            request, err := createRequestObject(r, user)
            if err != nil {
                logger.Error("Cannot create request object: " + err.Error())
                return
            }
            
            var query_filters string

            // logger.Info("Collection rules: " + collection.QueryRulesDirectoryPath + SEARCH_RULES_FILENAME)
            file_path := "access_rules/rules/" + strconv.FormatInt(collection.ID, 10) + "/" + SEARCH_RULES_FILENAME
			accept, err := ar.CheckRules(file_path, &query_filters, request, nil)
			if err != nil {
				logger.Error("Cannot check query rules: " + err.Error())
				http.Error(w, "An error occured", http.StatusInternalServerError)
				return
			}

			if accept {
                logger.Info("Access rules accepted")
                // analizing the filter
                used_collections, isValid, tokens := query_lang_anayzer.AnalyzeUserFilter(msg.CollectionName, msg.Filter)
                analyzed_expand := make([]querylang_types.TokenValueVariable, len(msg.Expand))
                for i, exp := range msg.Expand {
                    var token querylang_types.TokenValueVariable
                    token.Parts = strings.Split(exp, ".")
                    valid, exp_used_collections := query_lang_anayzer.AnalyzeVariableParts(msg.CollectionName, &token, querylang_types.ANALYZE_VARIABLES_PARTS_ANALYZE_TYPE_JOIN)
                    if !valid {
                        logger.Info("Expand is not valid")
                        return
                    }
                    
                    used_collections.Union(exp_used_collections)
                    analyzed_expand[i] = token
                }
                if isValid {
                    logger.Info("Filter is valid")
                    used_collections_filters := make(map[string]string)
                    used_collections_filters[msg.CollectionName] = query_filters
                    // checking used collections access rules
                    for collection_name, used_collection := range used_collections.GetMap() {
                        file_path := "access_rules/rules/" + strconv.FormatInt(used_collection.ID, 10) + "/" + SEARCH_RULES_FILENAME
                        var filters string
                        accept, err := ar.CheckRules(file_path, &filters, request, nil)
                        if err != nil {
                            logger.Error("Cannot check query rules: ", err)
                            http.Error(w, "An error occured", http.StatusInternalServerError)
                            return
                        }
                        if !accept {
                            http.Error(w, "Do not have access to collection " + used_collection.Name, http.StatusNotFound)
                            return
                        }

                        used_collections_filters[collection_name] = filters
                    }

                    // running query
                    resJson, err := db.RunUserCustomQuery(msg.CollectionName, msg.Fields, tokens, msg.SortBy, analyzed_expand, msg_limit, msg_offset, used_collections_filters)

                    if err != nil {
                        logger.Error("Cannot run query: " + err.Error())
                        http.Error(w, "An error occured", http.StatusInternalServerError)
                        return
                    }

                    
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusOK)
                    w.Write(resJson)
                }
			} else {
                logger.Warn("Access rules rejected")
                http.Error(w, "Do not have access to collection " + collection.Name, http.StatusNotFound)
                return
            }
        },
    )
}
