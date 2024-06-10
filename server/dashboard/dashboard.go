package dashboard

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http" // Used for build HTTP servers and clients.
	"os"

	"github.com/shachar1236/Baasa/database"
)

const ADMIN_ID = 1

// Port we listen on.
const portNum string = ":8080"

func RunDashboard(ctx context.Context, err_channel chan error) {
    logFile, err := os.OpenFile("logs/dashboard.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        err_channel <- err
    }
    mw := io.MultiWriter(os.Stdout, logFile)
    logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))
    // checking if admin exists

    var admin_exists bool
    var admin_session string

    admin_exists, err = database.DoesUserExistsById(ctx,ADMIN_ID)
    if err != nil {
        logger.Error("Got error while trying to check if admin exists: ", err)
        admin_exists = false
    }

    if admin_exists {
        admin, err := database.GetUserById(ctx, ADMIN_ID)
        if err != nil {
            err_channel <- errors.New("Admin should exists but cant find it")
        }
        admin_session = admin.Session
    }

    mux := http.NewServeMux()

    addRoutes(mux, logger, &admin_exists, &admin_session)

    logger.Info("Started on port " + portNum)
    logger.Info("To close connection CTRL+C :-)")

    // Spinning up the server.
    err = http.ListenAndServe(portNum, mux)
    if err != nil {
        err_channel <- err
    }
}
