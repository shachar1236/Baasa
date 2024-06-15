package dashboard

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/shachar1236/Baasa/database"
	"github.com/shachar1236/Baasa/database/types"
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

    mux.Handle("GET /collection", adminOnly(handleCollection(logger)))
    mux.Handle("GET /GetCollections", adminOnly(handleGetCollections(logger, db)))
    mux.Handle("GET /GetCollection", adminOnly(handleGetCollection(logger, db)))
    mux.Handle("POST /AddCollection", adminOnly(handleAddCollection(logger, db)))
    mux.Handle("DELETE /DeleteCollection", adminOnly(handleDeleteCollection(logger, db)))
    mux.Handle("POST /SaveCollectionChanges", adminOnly(handleSaveCollectionChanges(logger, db)))

    mux.Handle("POST /SetById", adminOnly(handleSetById(logger, db)))
    mux.Handle("DELETE /DeleteById", adminOnly(handleDeleteById(logger, db)))
    mux.Handle("POST /AddWithArgs", adminOnly(handleAddWithArgs(logger, db)))
    mux.Handle("/GetData", adminOnly(handleGetData(logger, db)))
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

func handleCollection(logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/collection.html")
        },
    )
}

func handleGetCollections(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            collections, err := db.GetCollections(r.Context())
            if err != nil {
                logger.Error("Got error while trying to retrieve collections: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            
            err = encode(w, r, http.StatusOK, collections)
            if err != nil {
                logger.Error("Got error while trying to encode collections: ", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        },
    )
}

func handleGetCollection(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            err := r.ParseForm()
            if err != nil {
                logger.Error("Got error while trying to parse form: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            collection_id, err := strconv.Atoi(r.Form.Get("id"))
            if err != nil {
                logger.Error("Got error while trying to retrieve collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            logger.Info("Got collection id: " + strconv.Itoa(collection_id))
            
            collection, err := db.GetCollectionById(r.Context(), int64(collection_id))
            err = encode(w, r, http.StatusOK, collection)
            if err != nil {
                logger.Error("Got error while trying to encode collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        },
    )
}

func handleAddCollection(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            collection, err := db.AddCollection(r.Context())
            if err != nil {
                logger.Error("Got error while trying to add collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            logger.Info("Created collection: " + collection.Name)

            err = encode(w, r, http.StatusOK, collection)
            if err != nil {
                logger.Error("Got error while trying to encode collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        },
    )
}

func handleDeleteCollection(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            q := r.URL.Query()
            name := q.Get("name")
            err := db.DeleteCollection(r.Context(), name)
            if err != nil {
                logger.Error("Got error while trying to delete collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        },
    )
}

func handleSaveCollectionChanges(logger *slog.Logger, db database.Database) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            // TODO: change to decodeValid to see if the changes the user has made are valid
            collection, err := decode[types.Collection](r);
            if err != nil {
                logger.Error("Got error while trying to decode collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            err = db.SaveCollectionChanges(r.Context(), collection)
            if err != nil {
                logger.Error("Got error while trying to save collection: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        },
    )
}

func handleSetById(logger *slog.Logger, db database.Database) http.HandlerFunc {
    type setMsg struct {
        CollectionName string `json:"collection_name"`
        ColumnName string `json:"column_name"`
        Id int64 `json:"id"`
        To any `json:"to"`
    }
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            params, err := decode[setMsg](r)
            if err != nil {
                logger.Error("Got error while trying to decode setMsg: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            err = db.ActionSetById(r.Context(), params.CollectionName, params.ColumnName, params.Id, params.To)
            if err != nil {
                logger.Error("Got error while trying to set by id: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }
            w.WriteHeader(http.StatusOK)
        },
    )
}

func handleGetData(logger *slog.Logger, db database.Database) http.HandlerFunc {
    type msg struct {
        Collection_name string `json:"collection_name"`
        From int64 `json:"from"`
        To int64 `json:"to"`
    }
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            params, err := decode[msg](r)
            if err != nil {
                logger.Error("Got error while trying to decode msg: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }

            data, err := db.ActionGetCollectionDataWithLimit(r.Context(), params.Collection_name, params.From, params.To)
            if err != nil {
                logger.Error("Got error while trying to get data: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }

            err = encode(w, r, http.StatusOK, data)
            if err != nil {
                logger.Error("Got error while trying to encode data: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }
        },
    )
}

func handleDeleteById(logger *slog.Logger, db database.Database) http.HandlerFunc {
    type deleteMsg struct {
        CollectionName string `json:"collection_name"`
        Id int64 `json:"id"`
    }
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            params, err := decode[deleteMsg](r)
            if err != nil {
                logger.Error("Got error while trying to decode deleteMsg: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            err = db.ActionDeleteById(r.Context(), params.CollectionName, params.Id)
            if err != nil {
                logger.Error("Got error while trying to delete by id: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }
            w.WriteHeader(http.StatusOK)
        },
    )
}

func handleAddWithArgs(logger *slog.Logger, db database.Database) http.HandlerFunc {
    type addMsg struct {
        CollectionName string `json:"collection_name"`
        Args map[string]any `json:"args"`
    }

    type retMsg struct {
        Id int64 `json:"id"`
    }
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

            params, err := decode[addMsg](r)
            if err != nil {
                logger.Error("Got error while trying to decode addMsg: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            
            id, err := db.ActionAdd(r.Context(), params.CollectionName, params.Args)
            if err != nil {
                logger.Error("Got error while trying to add with args: " + err.Error())
                w.WriteHeader(http.StatusBadRequest)
                encode(w, r, http.StatusBadRequest, err)
                return
            }
            encode(w, r, http.StatusOK, retMsg{
                Id: id,
            })
            w.WriteHeader(http.StatusOK)
        },
    )
}

