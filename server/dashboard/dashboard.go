package dashboard

import (
	"context"
	"log"      // logging messages to the console.
	"net/http" // Used for build HTTP servers and clients.

	"github.com/shachar1236/Baasa/database"
)

const ADMIN_ID = 1

var admin_exists bool
var admin_session string

// Port we listen on.
const portNum string = ":8080"

func IsAuthenticated(r *http.Request) bool {
    for _, cookie := range r.Cookies() {
        if cookie.Name == "session" {
            if cookie.Value == admin_session {
                return true
            } else {
                return false
            }
        }
    }
    return false
}

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
    if admin_exists {
        if IsAuthenticated(r) {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/home.html")
        } else {
            http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/login.html")
        }
    } else {
        http.ServeFile(w, r, "./dashboard/webpages/shachar_base/dist/register.html")
    }
}

func Register(w http.ResponseWriter, r *http.Request) {
    log.Println("Registering...")
    if !admin_exists && !IsAuthenticated(r) {
        log.Println("Admin dosent exist and user not authoticated")
        err := r.ParseForm()
        if err != nil {
            log.Println("Got error while trying to parse form: ", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        username := r.Form.Get("username")
        password := r.Form.Get("password")
        user_session, err := database.CreateUser(r.Context(), username, password)
        if err != nil {
            log.Println("Cant create user: ", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        log.Println("Created user")
        http.SetCookie(w, &http.Cookie{Name: "session", Value: user_session})
        admin_exists = true
        admin_session = user_session
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }

    w.WriteHeader(http.StatusBadRequest)
}

func Login(w http.ResponseWriter, r *http.Request) {
    log.Println("Logging...")
    if admin_exists && !IsAuthenticated(r) {
        log.Println("Admin exists and user is not authenticated")

        err := r.ParseForm()
        if err != nil {
            log.Println("Got error while trying to parse form: ", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        username := r.Form.Get("username")
        password := r.Form.Get("password")
        user, err := database.GetUser(r.Context(), username, password)
        if err != nil {
            log.Println("Got error while trying to retrieve user: ", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if user.Username == username && user.ID == ADMIN_ID {
            log.Println("Username and password good, logging in")
            http.SetCookie(w, &http.Cookie{Name: "session", Value: user.Session})
            http.Redirect(w, r, "/", http.StatusSeeOther)
        } else {
            w.WriteHeader(http.StatusBadRequest)
        }
    } else {
        http.Redirect(w, r, "/", http.StatusBadRequest)
    }
}

func GetTables(w http.ResponseWriter, r *http.Request) {
    if admin_exists && IsAuthenticated(r) {
        GetTables
    } else {
        w.WriteHeader(http.StatusBadRequest)
    }
}

func RunDashboard(ctx context.Context) {
    // checking if admin exists
    var err error
    admin_exists, err = database.DoesUserExistsById(ctx,ADMIN_ID)
    if err != nil {
        log.Println("Got error while trying to check if admin exists: ", err)
        admin_exists = false
    }

    if admin_exists {
        admin, err := database.GetUserById(ctx, ADMIN_ID)
        if err != nil {
            panic("Admin should exists but cant find it")
        }
        admin_session = admin.Session
    }

    fs := http.FileServer(http.Dir("./dashboard/webpages/shachar_base/dist/assets/"))
    mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

    // Registering our handler functions, and creating paths.
    mux.HandleFunc("/", Home)
    mux.HandleFunc("POST /register", Register)
    mux.HandleFunc("POST /login", Login)

    log.Println("Started on port", portNum)
    log.Println("To close connection CTRL+C :-)")

    // Spinning up the server.
    err = http.ListenAndServe(portNum, mux)
    if err != nil {
        log.Fatal(err)
    }
}
