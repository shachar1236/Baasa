package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
	"github.com/shachar1236/Baasa/database/types"
)

const SEARCH_RULES_FILENAME = "search.lua"

func addRoutes(mux *http.ServeMux, logger *slog.Logger, db database.Database, access_rules *access_rules.AccessRules) {
	mux.Handle("/query", handleQuery(logger, db, access_rules))
	mux.Handle("/collection/search", handleCollectionSearch(logger, db, access_rules))
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

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(resJson)
			}
		},
	)
}

func handleCollectionSearch(logger *slog.Logger, db database.Database, ar *access_rules.AccessRules) http.Handler {
    type SearchMessage struct {
        CollectionName string `json:"CollectionName"`
        Session string `json:"Session"`
        Fields []string `json:"Fields"`
        Filter string `json:"Filter"`
        SortBy string `json:"SortBy"`
        Expand string `json:"Expand"`
    }

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            var err error
            msg, err := decode[SearchMessage](r)
			if err != nil {
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
                used_collections, isValid, tokens := ar.AnalyzeUserFilter(msg.CollectionName, msg.Filter)
                if isValid {
                    logger.Info("Filter is valid")
                    used_collections_filters := make(map[string]string)
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
                    db.RunUserCustomQuery(msg.CollectionName, msg.Fields, tokens, msg.SortBy, msg.Expand, used_collections_filters)
                }
			}
        },
    )
}
