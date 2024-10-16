package usecase_test

import (
	"context"
	"testing"

	"github.com/appxpy/hive-test/internal/entity"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	ctx  context.Context

	// Test variables
	somePassword string
	someUser     *entity.User

	// Mocked units
	mockUserRepo *MockUserRepo

	// Tested usecase
	userUseCase usecase.UserUseCase

	// JWT Secret
	jwtSecret string
}

func (t *UserUseCaseSuite) SetupSuite() {
	t.somePassword = "kapusta"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(t.somePassword), bcrypt.DefaultCost)
	t.Require().NoError(err, "bcrypt hash error")

	t.someUser = &entity.User{
		Username:     "aboba",
		PasswordHash: string(passwordHash),
	}

	t.jwtSecret = "secretkey"
}

func (t *UserUseCaseSuite) SetupTest() {
	t.ctx = context.Background()
	t.ctrl = gomock.NewController(t.T())
	t.mockUserRepo = NewMockUserRepo(t.ctrl)
	t.userUseCase = usecase.NewUserUseCase(t.mockUserRepo, t.jwtSecret)
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

func (t *UserUseCaseSuite) TestRegister_GreenPath() {
	t.mockUserRepo.EXPECT().CreateUser(t.ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, user *entity.User) error {
			t.Equal(t.someUser.Username, user.Username)

			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(t.somePassword))
			t.NoError(err)

			return nil
		},
	)

	err := t.userUseCase.Register(t.ctx, t.someUser.Username, t.somePassword)

	t.NoError(err)
}

func (t *UserUseCaseSuite) TestRegister_ReturnsError_WhenRepoReturnsError() {
	t.mockUserRepo.EXPECT().CreateUser(t.ctx, gomock.Any()).Return(assert.AnError)

	err := t.userUseCase.Register(t.ctx, t.someUser.Username, t.somePassword)

	t.ErrorIs(err, assert.AnError)
}

func (t *UserUseCaseSuite) TestLogin_GreenPath() {
	t.mockUserRepo.EXPECT().GetUserByUsername(t.ctx, t.someUser.Username).Return(t.someUser, nil)

	token, err := t.userUseCase.Login(t.ctx, t.someUser.Username, t.somePassword)

	t.NoError(err)
	t.NotEmpty(token)

	// Verify JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.jwtSecret), nil
	})
	t.NoError(err)
	t.True(parsedToken.Valid)
	claims := parsedToken.Claims.(jwt.MapClaims)
	t.Equal(float64(t.someUser.ID), claims["user_id"])
}

func (t *UserUseCaseSuite) TestLogin_ReturnsError_WhenPasswordIncorrect() {
	t.mockUserRepo.EXPECT().GetUserByUsername(t.ctx, t.someUser.Username).Return(t.someUser, nil)

	token, err := t.userUseCase.Login(t.ctx, t.someUser.Username, "wrongpassword")

	t.Error(err)
	t.Contains(err.Error(), "invalid credentials")
	t.Empty(token)
}

func (t *UserUseCaseSuite) TestLogin_ReturnsError_WhenUserNotFound() {
	t.mockUserRepo.EXPECT().GetUserByUsername(t.ctx, t.someUser.Username).Return(nil, assert.AnError)

	token, err := t.userUseCase.Login(t.ctx, t.someUser.Username, t.somePassword)

	t.ErrorIs(err, assert.AnError)
	t.Empty(token)
}
