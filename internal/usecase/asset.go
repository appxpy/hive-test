package usecase

import (
	"context"
	"errors"

	"github.com/appxpy/hive-test/internal/entity"
)

// AssetUseCaseImpl implements the AssetUseCase interface.
type AssetUseCaseImpl struct {
	repo AssetRepo
}

// NewAssetUseCase creates a new AssetUseCase.
func NewAssetUseCase(repo AssetRepo) AssetUseCase {
	return &AssetUseCaseImpl{
		repo: repo,
	}
}

// AddAsset adds a new asset for a user.
func (uc *AssetUseCaseImpl) AddAsset(ctx context.Context, asset *entity.Asset) error {
	return uc.repo.CreateAsset(ctx, asset)
}

// RemoveAsset removes an asset owned by the user.
func (uc *AssetUseCaseImpl) RemoveAsset(ctx context.Context, assetID, userID int64) error {
	return uc.repo.DeleteAsset(ctx, assetID, userID)
}

// PurchaseAsset allows a user to purchase an asset.
func (uc *AssetUseCaseImpl) PurchaseAsset(ctx context.Context, assetID, buyerID int64) error {
	return uc.repo.ExecuteTx(ctx, func(repo AssetRepo) error {
		asset, err := repo.GetAssetByID(ctx, assetID, true)
		if err != nil {
			return err
		}

		if asset == nil {
			return errors.New("asset not found")
		}

		if asset.UserID == buyerID {
			return errors.New("cannot purchase your own asset")
		}

		err = repo.UpdateAssetOwner(ctx, assetID, buyerID)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetAssetsByUser retrieves all assets owned by the user.
func (uc *AssetUseCaseImpl) GetAssetsByUser(ctx context.Context, userID int64) ([]*entity.Asset, error) {
	return uc.repo.GetAssetsByUserID(ctx, userID)
}
