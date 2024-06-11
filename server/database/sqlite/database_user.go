package sqlite

import (
	"context"
	"crypto/sha256"

	"github.com/shachar1236/Baasa/database/sqlite/objects"
	"github.com/shachar1236/Baasa/database/types"
	"github.com/shachar1236/Baasa/database/utils"
)


func (this *SqliteDB) DoesUserExists(ctx context.Context, username string, password_hash types.PasswordHash) (bool, error) {
    count, err := this.objects_queries.CountUsersWithNameAndPassword(ctx, 
    objects.CountUsersWithNameAndPasswordParams{ Username: username, PasswordHash: password_hash})
    return count > 0, err
}

func (this *SqliteDB) DoesUserExistsById(ctx context.Context, id int64) (bool, error) {
    count, err := this.objects_queries.CountUsersWithId(ctx, id)
    return count > 0, err
}

func (this *SqliteDB) GetUserById(ctx context.Context, id int64) (types.User, error) {
    user, err := this.objects_queries.GetUserById(ctx, id)
    if err != nil {
        this.logger.Error("Cannot get user: ", err)
        return types.User{}, err
    }
    pass_hash := user.PasswordHash.([]byte)

    ret := types.User{
        ID: user.ID,
        Username: user.Username,
        Session: user.Session,
        PasswordHash: types.PasswordHash(pass_hash),
    }
    return ret, err
}

func (this *SqliteDB) GetUserBySession(ctx context.Context, session string) (types.User, error) {
    user, err := this.objects_queries.GetUserBySession(ctx, session)
    if err != nil {
        this.logger.Error("Cannot get user: ", err)
        return types.User{}, err
    }
    pass_hash := user.PasswordHash.([]byte)

    ret := types.User{
        ID: user.ID,
        Username: user.Username,
        Session: user.Session,
        PasswordHash: types.PasswordHash(pass_hash),
    }
    return ret, err
}

func (this *SqliteDB) CreateUser(ctx context.Context, username string, password string) (session string, err error){
    // generating password hash
    password_hash := sha256.Sum256([]byte(password))
    // generate random session
    rand_session, err := utils.GenerateSecureToken()
    if err != nil {
        this.logger.Error("Cannot generate session: ", err)
        return "", err
    }
    // creating user
    err = this.objects_queries.CreateUser(ctx, 
        objects.CreateUserParams{Username: username, PasswordHash: password_hash[:], Session: rand_session})
    return rand_session, err
}

func (this *SqliteDB) GetUser(ctx context.Context, username string, password string) (types.User, error) {
    password_hash := sha256.Sum256([]byte(password))
    user, err := this.objects_queries.GetUserByNameAndPassword(ctx, 
        objects.GetUserByNameAndPasswordParams{Username: username, PasswordHash: password_hash[:]}) 
    if err != nil {
        this.logger.Error("Cannot get user: ", err)
        return types.User{}, err
    }
    pass_hash := user.PasswordHash.([]byte)

    ret := types.User{
        ID: user.ID,
        Username: user.Username,
        Session: user.Session,
        PasswordHash: types.PasswordHash(pass_hash),
    }
    return ret, err
}

