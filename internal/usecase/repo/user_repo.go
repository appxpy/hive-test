package repo

import (
	"context"

	"github.com/appxpy/hive-test/internal/entity"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type UserRepoImpl struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) usecase.UserRepo {
	return &UserRepoImpl{DB: db}
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, user.Username, user.PasswordHash)
	return err
}

func (r *UserRepoImpl) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := r.DB.GetContext(ctx, user, query, username)
	return user, err
}
