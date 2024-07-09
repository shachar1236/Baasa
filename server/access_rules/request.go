package access_rules

import (
	"github.com/shachar1236/Baasa/database/types"
)

type Request struct {
    Method string
    Headers map[string]string
    Data map[string]string
    Auth types.User
}
