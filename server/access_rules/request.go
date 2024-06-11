package access_rules

import (
	"net/http"

	"github.com/shachar1236/Baasa/database/types"
)

type Request struct {
    Method string
    Headers http.Header
    Auth types.User
}
