package api

import (
	"context"
	"io"
	"os"
	"log/slog"
	"net/http"
)

const portNum string = ":5050"


func RunApi(ctx context.Context, err_channel chan error) {
	mux := http.NewServeMux()

    logFile, err := os.OpenFile("logs/api.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        err_channel <- err
    }
    mw := io.MultiWriter(os.Stdout, logFile)
    logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))

    addRoutes(mux, logger)

	logger.Info("Started on port " + portNum)
	logger.Info("To close connection CTRL+C :-)")

    err = http.ListenAndServe(portNum, mux)
	if err != nil {
        err_channel <- err
    }
}
