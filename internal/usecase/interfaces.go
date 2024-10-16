// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/appxpy/hive-test/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

// UserUseCase defines methods related to user operations.
type UserUseCase interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error)
}

// UserRepo defines methods to interact with the users in the database.
type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

// AssetUseCase defines methods related to asset operations.
type AssetUseCase interface {
	AddAsset(ctx context.Context, asset *entity.Asset) error
	RemoveAsset(ctx context.Context, assetID, userID int64) error
	PurchaseAsset(ctx context.Context, assetID, buyerID int64) error
	GetAssetsByUser(ctx context.Context, userID int64) ([]*entity.Asset, error)
}

// AssetRepo defines methods to interact with assets in the database.
type AssetRepo interface {
	CreateAsset(ctx context.Context, asset *entity.Asset) error
	DeleteAsset(ctx context.Context, assetID, userID int64) error
	GetAssetByID(ctx context.Context, assetID int64, forUpdate bool) (*entity.Asset, error)
	GetAssetsByUserID(ctx context.Context, userID int64) ([]*entity.Asset, error)
	UpdateAssetOwner(ctx context.Context, assetID, newOwnerID int64) error
	ExecuteTx(ctx context.Context, fn func(repo AssetRepo) error) error
}
