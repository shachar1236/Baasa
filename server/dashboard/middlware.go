package dashboard

import (
	"net/http"
)

// middleware is a function that wraps http.Handlers
// proving functionality before and after execution
// of the h handler.
type middleware func(h http.HandlerFunc) http.HandlerFunc

func newAdminOnlyMiddleware(admin_exists *bool, admin_session *string) middleware {
    return func(h http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if *admin_exists {
                for _, cookie := range r.Cookies() {
                    if cookie.Name == "session" {
                        if cookie.Value == *admin_session {
                            h(w, r)
                            return
                        } else {
                            http.NotFound(w, r)
                            return
                        }
                    }
                }
            } else {
                http.NotFound(w, r)
            }
            return
        })
    }
}

func newGuestOnlyMiddleware(admin_exists *bool, admin_session *string) middleware {
    return func(h http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if *admin_exists {
                for _, cookie := range r.Cookies() {
                    if cookie.Name == "session" {
                        if cookie.Value == *admin_session {
                            http.NotFound(w, r)
                            return
                        }                 }
                }
            }

            h(w, r)
            return
        })
    }
}

func newAdminDosentExistsMiddleware(admin_exists *bool) middleware {
    return func(h http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if *admin_exists {
                http.NotFound(w, r)
            }

            h(w, r)
            return
        })
    }
}
