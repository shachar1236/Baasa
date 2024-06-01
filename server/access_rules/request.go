package access_rules

import (
	"net/http"

	"github.com/shachar1236/Baasa/database/objects"
)

type Request struct {
    Method string
    Headers http.Header
    Auth objects.User
}
