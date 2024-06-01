package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
)

const portNum string = ":5050"

func query(w http.ResponseWriter, r *http.Request) {
    var msg QueryMessage 

    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    log.Println(msg.QueryId, msg.Session)

    query_filters := ""

    // getting query
    query, err := database.GetQuaryById(r.Context(), msg.QueryId)
    if err != nil {
        log.Println("Cannot get query: ", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // getting user
    user, err := database.GetUserBySession(r.Context(), msg.Session)
    if err != nil {
        log.Println("Cannot get user")
        user.ID = -1
        user.Username = ""
    }

    // check if query is by the rules
    request := access_rules.Request{
        Method: r.Method,
        Headers: r.Header,
        Auth: user,
    }


    accept, err := access_rules.CheckRules(query.QueryRulesFilePath, &query_filters, request)
    if err != nil {
        log.Println("Cannot check query rules: ", err)
        http.Error(w, "An error occured", http.StatusInternalServerError)
        return
    }

    if accept {
        // run query
        resJson, err := database.RunQueryWithFilters(r.Context(), query, msg.QueryArgs, query_filters)
        if err != nil {
            log.Println("Cannot run query: ", err)
            http.Error(w, "An error occured", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(resJson))
    }

}

func RunApi(ctx context.Context) {
    mux := http.NewServeMux()
    mux.HandleFunc("/query", query)
    log.Println("Started on port", portNum)
    log.Println("To close connection CTRL+C :-)")
    log.Fatal(http.ListenAndServe(portNum, mux))
}
