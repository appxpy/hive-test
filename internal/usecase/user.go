package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/appxpy/hive-test/internal/entity"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

// UserUseCaseImpl implements UserUseCase.
type UserUseCaseImpl struct {
	Repo      UserRepo
	JWTSecret string
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(repo UserRepo, jwtSecret string) UserUseCase {
	return &UserUseCaseImpl{
		Repo:      repo,
		JWTSecret: jwtSecret,
	}
}

// Register registers a new user.
func (uc *UserUseCaseImpl) Register(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return uc.Repo.CreateUser(ctx, user)
}

// Login authenticates a user and returns a JWT token.
func (uc *UserUseCaseImpl) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.Repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
