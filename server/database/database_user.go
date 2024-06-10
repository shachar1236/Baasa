package database

import (
	"context"
	"crypto/sha256"

	"github.com/shachar1236/Baasa/database/objects"
)


func DoesUserExists(ctx context.Context, username string, password_hash PasswordHash) (bool, error) {
    count, err := objects_queries.CountUsersWithNameAndPassword(ctx, 
        objects.CountUsersWithNameAndPasswordParams{ Username: username, PasswordHash: password_hash[:]})
    return count > 0, err
}

func DoesUserExistsById(ctx context.Context, id int64) (bool, error) {
    count, err := objects_queries.CountUsersWithId(ctx, id)
    return count > 0, err
}

func GetUserById(ctx context.Context, id int64) (objects.User, error) {
    user, err := objects_queries.GetUserById(ctx, id)
    return user, err
}

func GetUserBySession(ctx context.Context, session string) (objects.User, error) {
    user, err := objects_queries.GetUserBySession(ctx, session)
    return user, err
}

func CreateUser(ctx context.Context, username string, password string) (session string, err error){
    // generating password hash
    password_hash := sha256.Sum256([]byte(password))
    // generate random session
    rand_session, err := GenerateSecureToken()
    if err != nil {
        return "", err
    }
    // creating user
    err = objects_queries.CreateUser(ctx, 
        objects.CreateUserParams{Username: username, PasswordHash: password_hash[:], Session: rand_session})
    return rand_session, err
}

func GetUser(ctx context.Context, username string, password string) (objects.User, error) {
    password_hash := sha256.Sum256([]byte(password))
    user, err := objects_queries.GetUserByNameAndPassword(ctx, 
        objects.GetUserByNameAndPasswordParams{Username: username, PasswordHash: password_hash[:]}) 
    return user, err
}

