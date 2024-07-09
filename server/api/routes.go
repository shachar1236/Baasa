package api

import (
	"log/slog"
	"net/http"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
)

func addRoutes(mux *http.ServeMux, logger *slog.Logger, db database.Database, access_rules *access_rules.AccessRules) {
	mux.Handle("/query", handleQuery(logger, db, access_rules))
	mux.Handle("/", http.NotFoundHandler())
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
			request := access_rules.Request{
				Method:  r.Method,
				Auth:    user,
			}
            for k, v := range r.Header {
                if len(v) > 0 {
                    request.Headers[k] = v[0]
                }
            }
            err = r.ParseForm()
            if err != nil {
                logger.Error("Cannot parse form: ", err)
            }
            for k, v := range r.Form {
                if len(v) > 0 {
                    request.Data[k] = v[0]
                }
            }

			accept, err := ar.CheckRules(query.QueryRulesFilePath, &query_filters, request)
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
    type QueryMessage struct {
        QueryId int64 `json:"QueryId"`
        Session string `json:"Session"`
        QueryArgs map[string]any `json:"QuaryArgs"`
    }

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            var err error
            // msg, err := decode[QueryMessage](r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
        },
    )
}
