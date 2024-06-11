package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
)

const portNum string = ":5050"


func RunApi(ctx context.Context, err_channel chan error, db database.Database, access_rules access_rules.AccessRules) {
	mux := http.NewServeMux()

    logFile, err := os.OpenFile("logs/api.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    var mw io.Writer
    if err != nil {
        mw = os.Stdout
    } else {
        mw = io.MultiWriter(os.Stdout, logFile)
    }
    logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))

    addRoutes(mux, logger, db, access_rules)

	logger.Info("Started on port " + portNum)
	logger.Info("To close connection CTRL+C :-)")

    err = http.ListenAndServe(portNum, mux)
	if err != nil {
        err_channel <- err
    }
}
