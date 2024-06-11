package dashboard

import (
	"log/slog"
	"net/http"

	"github.com/shachar1236/Baasa/database"
)

func addRoutes(mux *http.ServeMux, logger *slog.Logger, db database.Database, admin_exists *bool, admin_session *string) {
    fs := http.FileServer(http.Dir("./dashboard/webpages/shachar_base/dist/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

    adminOnly := newAdminOnlyMiddleware(admin_exists, admin_session)
    guestOnly := newGuestOnlyMiddleware(admin_exists, admin_session)
    adminDosentExists := newAdminDosentExistsMiddleware(admin_exists)

    // Registering our handler functions, and creating paths.
    mux.Handle("/", adminOnly(handleHome(logger)))
    mux.Handle("POST /register", adminDosentExists(handleRegisterPost(logger, db, admin_exists, admin_session)))
    mux.Handle("POST /login", guestOnly(handleLoginPost(logger, db)))

    mux.Handle("GET /register", adminDosentExists(handleRegisterGet(logger)))
    mux.Handle("GET /login", guestOnly(handleLoginGet(logger)))
}

func handleHome(logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/home.html")
            w.WriteHeader(http.StatusOK)
		},
	)
}

func handleRegisterPost(logger *slog.Logger, db database.Database, admin_exists *bool, admin_session *string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            logger.Info("Registering...")
            err := r.ParseForm()
            if err != nil {
                logger.Error("Got error while trying to parse form: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            username := r.Form.Get("username")
            password := r.Form.Get("password")
            user_session, err := db.CreateUser(r.Context(), username, password)
            if err != nil {
                logger.Error("Cant create user: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            logger.Info("Created user")
            http.SetCookie(w, &http.Cookie{Name: "session", Value: user_session})
            *admin_exists = true
            *admin_session = user_session
            http.Redirect(w, r, "/", http.StatusSeeOther)
		},
	)
}

func handleRegisterGet(logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/register.html")
        },
    )
}

func handleLoginPost(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            logger.Info("Logging...")
            logger.Info("Admin exists and user is not authenticated")

            err := r.ParseForm()
            if err != nil {
                logger.Error("Got error while trying to parse form: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            username := r.Form.Get("username")
            password := r.Form.Get("password")
            user, err := db.GetUser(r.Context(), username, password)

            if err != nil {
                logger.Error("Got error while trying to retrieve user: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            if user.Username == username && user.ID == ADMIN_ID {
                logger.Info("Username and password good, logging in")
                http.SetCookie(w, &http.Cookie{Name: "session", Value: user.Session})
                http.Redirect(w, r, "/", http.StatusSeeOther)
            } else {
                w.WriteHeader(http.StatusBadRequest)
            }
		},
	)
}

func handleLoginGet(logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/login.html")
        },
    )
}
